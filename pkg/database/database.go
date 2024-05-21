package database

import (
	"context"
	"fmt"
	"merge-api/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	DB *pgxpool.Pool
}

func NewDatabaseConnection(config *config.Config) (Database, error) {
	dsn := makeDatabaseString(config)
	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return Database{}, err
	}

	pgxConfig.MaxConns = 20
	pgxConfig.MinConns = 1
	pgxConfig.MaxConnIdleTime = 5 * time.Minute
	pgxConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return Database{}, err
	}

	return Database{
		DB: pool,
	}, nil
}

func makeDatabaseString(config *config.Config) string {
	template := "postgresql://%v:%v@%v:%v/%v?sslmode=disable"
	dsn := fmt.Sprintf(template, config.Postgresql.User, config.Postgresql.Password,
		config.Postgresql.Address, config.Postgresql.Port, config.Postgresql.Database)
	return dsn
}
