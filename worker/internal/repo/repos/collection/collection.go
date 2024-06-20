package collection

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"merge-api/shared/pkg/database"
	"merge-api/worker/pkg"
)

var NoNextElementError = errors.New("no next element for such items")

type Repo struct {
	database *database.Database
}

func (r *Repo) GetNextCollectionItem(ctx context.Context, item pkg.CollectionItem) (pkg.CollectionItem, error) {
	conn, err := r.database.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	collectionId, collectionItemId := item.GetCollectionInfo()

	stmt := sq.
		Select("collection_id, level").
		From("collection_items").
		Where(sq.Or{
			sq.Eq{
				"collection_id": collectionId,
				"level":         collectionItemId + 1,
			},
			sq.Eq{
				"collection_id": collectionId + 1,
				"level":         1, // first item in collection
			},
		}).
		OrderBy("collection_id").
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return nil, err
		} else {
			err = rows.Scan(&collectionId, &collectionItemId)
			result := pkg.NewCollectionItemImpl(collectionId, collectionItemId)
			return &result, nil
		}
	}
	return nil, NoNextElementError
}

func (r *Repo) IsItemMergeable(ctx context.Context, item pkg.CollectionItem) (bool, error) {
	conn, err := r.database.DB.Acquire(ctx)
	if err != nil {
		return false, err
	}

	defer conn.Release()

	collectionId, collectionItemId := item.GetCollectionInfo()

	stmt := sq.
		Select("mergeable").
		From("collection_items").
		Where(sq.Eq{
			"collection_id": collectionId,
			"level":         collectionItemId,
		}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := stmt.ToSql()
	if err != nil {
		return false, err
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return false, err
	}

	var mergeable bool
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return false, err
		} else {
			err = rows.Scan(&mergeable)
			return mergeable, nil
		}
	}
	return false, NoNextElementError
}

func NewCollectionRepo(database *database.Database) *Repo {
	return &Repo{
		database: database,
	}
}
