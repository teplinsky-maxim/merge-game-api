package board

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"merge-api/shared/pkg/board"
	"merge-api/shared/pkg/board/inmemory"
	"merge-api/shared/pkg/database"
	"merge-api/worker/internal/repo/repos/collection/redis_board"
	"merge-api/worker/pkg"
)

type PgxTxType string

const TxKey = PgxTxType("tx")

var CoordinatesOutOfBoundsError = errors.New("coordinates out of bounds error")
var NoFreeCellsError = errors.New("no free cells are in the board")

type Repo struct {
	database *database.Database
}

func (r *Repo) GetBoard(ctx context.Context, id uint) (board.Board[pkg.CollectionItem], error) {
	panic("implement me")
}

func (r *Repo) GetBoardByCoordinates(ctx context.Context, id, w, h uint) (pkg.CollectionItem, error) {
	tx, err := r.database.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Release()

	stmt := sq.
		Select("board_cells.collection_id", "board_cells.collection_item_id", "boards.width", "boards.height").
		From("board_cells").
		LeftJoin("boards ON boards.id = board_cells.board_id").
		Where(sq.Eq{
			"board_id": id,
			"cell_h":   h,
			"cell_w":   w,
		}).
		PlaceholderFormat(sq.Dollar)
	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var collectionId, collectionItemId, maxW, maxH uint
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return nil, err
		} else {
			// TODO: in case of oob make sure it is also checked in redis
			err = rows.Scan(&collectionId, &collectionItemId, &maxW, &maxH)
			if w > maxW || h > maxH {
				return nil, CoordinatesOutOfBoundsError
			}
			result := pkg.NewCollectionItemImpl(collectionId, collectionItemId)
			return &result, nil
		}
	}
	return nil, redis_board.BoardCellEmptyError

}

func (r *Repo) CreateBoard(ctx context.Context, w, h uint) (board.Board[pkg.CollectionItem], uint, error) {
	tx, err := r.database.DB.Begin(ctx)
	if err != nil {
		return nil, 0, err
	}

	stmt := sq.
		Insert("boards").
		Columns("width", "height").
		Values(w, h).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)
	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}

	result := inmemory.NewBoard[pkg.CollectionItem](w, h)
	var id uint
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return nil, 0, err
		} else {
			err = rows.Scan(&id)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, 0, err
	}
	return &result, id, nil
}

func (r *Repo) UpdateCell(ctx context.Context, id, w, h uint, collectionItem pkg.CollectionItem) error {
	var tx pgx.Tx
	var err error
	txWasCreatedHere := false
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if !ok {
		txWasCreatedHere = true
		conn, err := r.database.DB.Acquire(ctx)
		defer conn.Release()
		tx, err = conn.Begin(ctx)
		if err != nil {
			return err
		}
	}

	collectionId, collectionItemId := collectionItem.GetCollectionInfo()

	stmt := sq.
		Insert("board_cells").
		Columns("board_id", "cell_w", "cell_h", "collection_id", "collection_item_id").
		Values(id, w, h, collectionId, collectionItemId).
		Suffix("ON CONFLICT (board_id, cell_w, cell_h) DO UPDATE SET collection_id = $4, collection_item_id = $5," +
			" time_created = now()").
		PlaceholderFormat(sq.Dollar)

	query, args, err := stmt.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if txWasCreatedHere {
		return tx.Commit(ctx)
	}

	return err
}

func (r *Repo) ClearCell(ctx context.Context, id, w, h uint) error {
	var tx pgx.Tx
	var err error
	txWasCreatedHere := false
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if !ok {
		txWasCreatedHere = true
		conn, err := r.database.DB.Acquire(ctx)
		defer conn.Release()
		tx, err = conn.Begin(ctx)
		if err != nil {
			return err
		}
	}

	stmt := sq.
		Delete("board_cells").
		Where(sq.Eq{
			"board_id": id,
			"cell_w":   w,
			"cell_h":   h,
		}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := stmt.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if txWasCreatedHere {
		return tx.Commit(ctx)
	}

	return err
}

func (r *Repo) FindEmptyCell(ctx context.Context, id uint) (uint, uint, error) {
	conn, err := r.database.DB.Acquire(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer conn.Release()

	query := fmt.Sprintf(`
WITH board_dimensions AS (SELECT b.id AS board_id, b.width, b.height
                          FROM boards b
                          WHERE b.id = %d),
     possible_cells AS (SELECT generate_series(1, (SELECT width FROM board_dimensions))  AS cell_w,
                               generate_series(1, (SELECT height FROM board_dimensions)) AS cell_h,
                               (SELECT board_id FROM board_dimensions)                   AS board_id),
     all_possible_cells AS (SELECT pw.cell_w, ph.cell_h, pw.board_id
                            FROM possible_cells pw
                                     CROSS JOIN possible_cells ph
                            WHERE pw.cell_w <= (SELECT width FROM board_dimensions)
                              AND ph.cell_h <= (SELECT height FROM board_dimensions)),
     occupied_cells AS (SELECT cell_w, cell_h, board_id
                        FROM board_cells
                        WHERE board_id = %d)
SELECT apc.cell_w, apc.cell_h
FROM all_possible_cells apc
         LEFT JOIN occupied_cells oc
                   ON apc.cell_w = oc.cell_w
                       AND apc.cell_h = oc.cell_h
                       AND apc.board_id = oc.board_id
WHERE oc.cell_w IS NULL
  AND oc.cell_h IS NULL
LIMIT 1;
`, id, id)

	var w, h uint
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return 0, 0, err
	}

	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return 0, 0, err
		} else {
			err = rows.Scan(&w, &h)
			if err != nil {
				return 0, 0, err
			}
			return w, h, nil
		}
	}
	return 0, 0, NoFreeCellsError
}

func (r *Repo) DeleteBoard(ctx context.Context, id uint) error {
	panic("implement me")
}

func NewBoardRepo(database *database.Database) *Repo {
	return &Repo{
		database: database,
	}
}
