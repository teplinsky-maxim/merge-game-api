package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"merge-api/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

const TestsSchemaName = "tests"

type DatabaseTest Database
type Database struct {
	DB *pgxpool.Pool
}

func NewDatabaseConnection(config *config.Config) (Database, error) {
	return newDatabaseConnection(config, false)
}

func NewDatabaseTestConnection(config *config.Config) (DatabaseTest, error) {
	connection, err := newDatabaseConnection(config, true)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	tx, err := connection.DB.Begin(ctx)
	if err != nil {
		panic(err)
	}
	stmt := "CREATE SCHEMA IF NOT EXISTS " + TestsSchemaName
	if err != nil {
		panic(err)
	}
	_, err = tx.Exec(ctx, stmt)
	if err != nil {
		panic(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		panic(err)
	}
	err = RunMigrations(config, TestsSchemaName)
	if err != nil {
		panic(err)
	}
	return DatabaseTest(connection), nil
}

func newDatabaseConnection(config *config.Config, testConnection bool) (Database, error) {
	schema := "public"
	if testConnection {
		schema = TestsSchemaName
	}
	dsn := makeDatabaseString(config, schema)
	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return Database{}, err
	}

	pgxConfig.MaxConns = 20
	pgxConfig.MinConns = 1
	pgxConfig.MaxConnIdleTime = 5 * time.Minute
	pgxConfig.HealthCheckPeriod = 1 * time.Minute

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return Database{}, err
	}

	return Database{
		DB: pool,
	}, nil
}

func makeDatabaseString(config *config.Config, schema string) string {
	template := "postgresql://%v:%v@%v:%v/%v?sslmode=disable&search_path=%v"
	dsn := fmt.Sprintf(template, config.Postgresql.User, config.Postgresql.Password,
		config.Postgresql.Address, config.Postgresql.Port, config.Postgresql.Database, schema)
	return dsn
}

func (r *DatabaseTest) SetUp(tableName string) {
	stmt := fmt.Sprintf("TRUNCATE TABLE %v RESTART IDENTITY", tableName)
	_, err := r.DB.Exec(context.Background(), stmt)
	if err != nil {
		panic(err)
	}
}
