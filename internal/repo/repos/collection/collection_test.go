package collection

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"merge-api/config"
	"merge-api/internal/entity"
	"merge-api/pkg/database"
	"reflect"
	"testing"
)

func prepareDatabaseWithRepo() (error, database.DatabaseTest, *CollectionRepo) {
	conf, err := config.NewConfigWithDiscover(nil)
	if err != nil {
		panic(err)
	}
	connection, err := database.NewDatabaseTestConnection(conf)
	if err != nil {
		panic(err)
	}
	repo := NewCollectionRepo((*database.Database)(&connection))
	return err, connection, repo
}

func prepareCollections(tx pgx.Tx, ctx context.Context) ([]entity.Collection, error) {
	collections := []entity.Collection{{
		Name: "test-name-1",
	}, {
		Name: "test-name-2",
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
	return collections, err
}

func prepareCollectionItems(tx pgx.Tx, ctx context.Context, cId1, cId2 uint) ([]entity.CollectionItem, error) {
	items := []entity.CollectionItem{{
		CollectionId: cId1,
		Name:         "test-name-1",
		Level:        1,
		Mergeable:    true,
		CanCreate:    false,
	}, {
		CollectionId: cId1,
		Name:         "test-name-2",
		Level:        2,
		Mergeable:    true,
		CanCreate:    true,
	}, {
		CollectionId: cId2,
		Name:         "another-test-name",
		Level:        1,
		Mergeable:    true,
		CanCreate:    false,
	}}

	stmt := sq.Insert("collection_items").
		Columns("collection_id", "name", "level", "mergeable", "can_create").
		Values(items[0].CollectionId, items[0].Name, items[0].Level, items[0].Mergeable, items[0].CanCreate).
		Values(items[1].CollectionId, items[1].Name, items[1].Level, items[1].Mergeable, items[1].CanCreate).
		Values(items[2].CollectionId, items[2].Name, items[2].Level, items[2].Mergeable, items[2].CanCreate).
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
			items[i].ID = id
			i++
		}
	}
	return items, err
}

func prepareCreationRules(tx pgx.Tx, ctx context.Context, collectionItemId, generatesCollectionItemId uint) ([]entity.CreationRule, error) {
	items := []entity.CreationRule{{
		CollectionItemId:         collectionItemId,
		GenerateCollectionItemId: generatesCollectionItemId,
	}}

	stmt := sq.Insert("creation_rules").
		Columns("collection_item_id", "generate_collection_item_id").
		Values(items[0].CollectionItemId, items[0].GenerateCollectionItemId).
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
			items[i].ID = id
			i++
		}
	}
	return items, err
}

func TestCollectionRepo_GetCollections(t *testing.T) {
	err, connection, repo := prepareDatabaseWithRepo()
	connection.SetUp("collections")

	ctx := context.Background()
	tx, err := connection.DB.Begin(ctx)
	if err != nil {
		panic(err)
	}

	collections, err := prepareCollections(tx, ctx)

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
			got, err := repo.GetCollections(tt.args.ctx, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCollections() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCollections() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollectionRepo_GetCollection(t *testing.T) {
	err, connection, repo := prepareDatabaseWithRepo()
	connection.SetUp("collections")
	connection.SetUp("collection_items")
	connection.SetUp("creation_rules")

	ctx := context.Background()
	tx, err := connection.DB.Begin(ctx)
	if err != nil {
		panic(err)
	}

	collections, err := prepareCollections(tx, ctx)
	collectionItems, err := prepareCollectionItems(tx, ctx, collections[0].ID, collections[1].ID)
	_, err = prepareCreationRules(tx, ctx, collectionItems[1].ID, collectionItems[2].ID)

	err = tx.Commit(ctx)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name         string
		collectionId uint
		want         entity.CollectionWithItems
		wantErr      bool
	}{
		{
			name:         "ok",
			collectionId: collections[0].ID,
			want:         entity.CollectionWithItems{},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetCollection(ctx, tt.collectionId)
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

func TestCollectionRepo_CreateCollection(t *testing.T) {
	err, connection, repo := prepareDatabaseWithRepo()
	if err != nil {
		panic(err)
	}
	connection.SetUp("collections")
	ctx := context.Background()

	tests := []struct {
		name    string
		args    string
		want    entity.Collection
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.CreateCollection(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateCollection() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollectionRepo_CreateCollectionItems(t *testing.T) {
	type fields struct {
		database *database.Database
	}
	type args struct {
		ctx   context.Context
		items []entity.CollectionItem
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entity.CollectionItem
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CollectionRepo{
				database: tt.fields.database,
			}
			got, err := c.CreateCollectionItems(tt.args.ctx, tt.args.items)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCollectionItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateCollectionItems() got = %v, want %v", got, tt.want)
			}
		})
	}
}
