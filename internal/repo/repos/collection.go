package repos

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"merge-api/internal/entity"
	"merge-api/pkg/database"
)

type CollectionRepo struct {
	database *database.Database
}

func (c *CollectionRepo) GetCollection(ctx context.Context, offset, limit uint) ([]entity.Collection, error) {
	conn, err := c.database.DB.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	query, _, err := sq.
		Select("*").
		From("public.collections").
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

func NewCollectionRepo(database *database.Database) *CollectionRepo {
	return &CollectionRepo{
		database: database,
	}
}
