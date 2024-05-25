package main

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"merge-api/config"
	"merge-api/internal/entity"
	"merge-api/internal/repo/repos"
	"merge-api/pkg/database"
	"time"
)

func main() {
	dsn := "postgresql://user:password@localhost:5432/merge?sslmode=disable&search_path=tests"
	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}

	pgxConfig.MaxConns = 20
	pgxConfig.MinConns = 1
	pgxConfig.MaxConnIdleTime = 5 * time.Minute
	pgxConfig.HealthCheckPeriod = 1 * time.Minute

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	conn, err := pool.Acquire(ctx)

	if err != nil {
		panic(err)
	}

	query, _, err := sq.
		Select("*").
		From("collections").
		Limit(uint64(10)).
		Offset(uint64(1)).ToSql()
	if err != nil {
		panic(err)
	}
	rows, err := conn.Query(ctx, query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var result []entity.Collection
	var r entity.Collection
	for rows.Next() {
		err = rows.Scan(&r.ID, &r.Name)
		if err != nil {
			panic(err)
		}
		result = append(result, r)
	}
	fmt.Printf("%v", r)

	conf, err := config.NewConfigWithDiscover(nil)
	if err != nil {
		panic(err)
	}
	connection, err := database.NewDatabaseTestConnection(conf)
	if err != nil {
		panic(err)
	}
	repo := repos.NewCollectionRepo((*database.Database)(&connection))
	collections, err := repo.GetCollections(ctx, 1, 10)
	if err != nil {
		return
	}
	fmt.Printf("%v", collections)

}
