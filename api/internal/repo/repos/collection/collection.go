package collection

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"merge-api/api/internal/entity"
	"merge-api/shared/pkg/database"
)

var (
	CollectionDoesNotExistsError = errors.New("collection with this id does not exists")
)

type CollectionRepo struct {
	database *database.Database
}

func (c *CollectionRepo) CreateCollection(ctx context.Context, name string) (entity.Collection, error) {
	tx, err := c.database.DB.Begin(ctx)
	if err != nil {
		return entity.Collection{}, nil
	}

	stmt := sq.
		Insert("collections").Columns("name").
		Values(name).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)
	query, args, err := stmt.ToSql()
	if err != nil {
		return entity.Collection{}, nil
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return entity.Collection{}, nil
	}

	var result entity.Collection
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return entity.Collection{}, nil
		} else {
			err = rows.Scan(&result.ID)
			result.Name = name
		}
	}
	return result, err
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

func (c *CollectionRepo) CreateCollectionItems(ctx context.Context, items []entity.CollectionItem) ([]entity.CollectionItem, error) {
	tx, err := c.database.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}

	stmt := sq.
		Insert("collection_items").
		Columns("collection_id", "name", "level", "mergeable", "can_create").
		Suffix("RETURNING id, collection_id, name, level, mergeable, can_create").
		PlaceholderFormat(sq.Dollar)

	for _, item := range items {
		stmt = stmt.Values(item.CollectionId, item.Name, item.Level, item.Mergeable, item.CanCreate)
	}

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var createdItems []entity.CollectionItem
	for rows.Next() {
		var item entity.CollectionItem
		err = rows.Scan(&item.ID, &item.CollectionId, &item.Name, &item.Level, &item.Mergeable, &item.CanCreate)
		if err != nil {
			return nil, err
		}
		createdItems = append(createdItems, item)
	}

	return createdItems, nil
}

func NewCollectionRepo(database *database.Database) *CollectionRepo {
	return &CollectionRepo{
		database: database,
	}
}
