package repos

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"merge-api/config"
	"merge-api/internal/entity"
	"merge-api/pkg/database"
	"reflect"
	"testing"
)

func TestCollectionRepo_GetCollection(t *testing.T) {
	conf, err := config.NewConfigWithDiscover(nil)
	if err != nil {
		panic(err)
	}
	connection, err := database.NewDatabaseTestConnection(conf)
	if err != nil {
		panic(err)
	}
	connection.SetUp("collections")
	repo := NewCollectionRepo((*database.Database)(&connection))

	ctx := context.Background()
	tx, err := connection.DB.Begin(ctx)
	if err != nil {
		panic(err)
	}

	collections := []entity.Collection{{
		Name: "name1qweqweq",
	}, {
		Name: "name2qweqweqwe",
	}}

	stmt := sq.Insert("collections").
		Columns("name").
		Values(collections[0].Name).
		Values(collections[1].Name).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := stmt.ToSql()
	if err != nil {
		panic(err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		panic(err)
	}

	i := 0
	for rows.Next() {
		var id uint
		err = rows.Err()
		if err != nil {
			panic(err)
		} else {
			err = rows.Scan(&id)
			collections[i].ID = id
			i++
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		panic(err)
	}

	type testArgs struct {
		ctx    context.Context
		offset uint
		limit  uint
	}
	tests := []struct {
		name    string
		args    testArgs
		want    []entity.Collection
		wantErr bool
	}{
		{
			name: "ok",
			args: testArgs{
				ctx:    context.Background(),
				offset: 0,
				limit:  100,
			},
			want:    collections,
			wantErr: false,
		},
		{
			name: "limit works",
			args: testArgs{
				ctx:    context.Background(),
				offset: 0,
				limit:  0,
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetCollection(tt.args.ctx, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCollection() got = %v, want %v", got, tt.want)
			}
		})
	}
}
