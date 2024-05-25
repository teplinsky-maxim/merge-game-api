package repos

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"merge-api/internal/entity"
	"merge-api/pkg/database"
)

var (
	CollectionDoesNotExistsError = errors.New("collection with this id does not exists")
)

type CollectionRepo struct {
	database *database.Database
}

func (c *CollectionRepo) GetCollections(ctx context.Context, offset, limit uint) ([]entity.Collection, error) {
	conn, err := c.database.DB.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	query, _, err := sq.
		Select("*").
		From("collections").
		Limit(uint64(limit)).
		Offset(uint64(offset)).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []entity.Collection
	for rows.Next() {
		var r entity.Collection
		err = rows.Scan(&r.ID, &r.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

func (c *CollectionRepo) GetCollection(ctx context.Context, collectionId uint) (entity.CollectionWithItems, error) {
	conn, err := c.database.DB.Acquire(ctx)

	defer conn.Release()
	if err != nil {
		return entity.CollectionWithItems{}, nil
	}

	stmt := sq.Expr(`
SELECT c.id,
       c.name,
       json_agg(
               json_build_object(
                       'id', ci.id,
                       'name', ci.name,
                       'level', ci.level,
                       'generates', cr_details
               )
       ) as items
FROM collections c
         JOIN collection_items ci on c.id = ci.collection_id
         LEFT JOIN (SELECT cr.id,
                           cr.collection_item_id,
                           cr.generate_collection_item_id,
                           json_build_object(
                                   'collection_id', gci.collection_id,
                                   'level', gci.level
                           ) as cr_details
                    FROM creation_rules cr
                             LEFT JOIN collection_items gci on cr.generate_collection_item_id = gci.id) cr
                   on ci.id = cr.collection_item_id
WHERE c.id = $1
GROUP BY c.id;
`, collectionId)

	query, args, err := stmt.ToSql()
	if err != nil {
		return entity.CollectionWithItems{}, nil
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return entity.CollectionWithItems{}, nil
	}

	defer rows.Close()
	for rows.Next() {
		result := new(entity.CollectionWithItems)
		err = rows.Scan(&result.ID, &result.Name, &result.Items)
		if err != nil {
			return entity.CollectionWithItems{}, err
		}
		return *result, err
	}
	return entity.CollectionWithItems{}, CollectionDoesNotExistsError
}

func NewCollectionRepo(database *database.Database) *CollectionRepo {
	return &CollectionRepo{
		database: database,
	}
}
