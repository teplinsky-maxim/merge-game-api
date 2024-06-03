package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"merge-api/shared/config"
	"os"
	"path/filepath"
	"strings"
)

func RunMigrations(config *config.Config, schema string) error {
	dsn := makeDatabaseString(config, schema)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	driver, err := pgx.WithInstance(db, &pgx.Config{
		SchemaName: schema,
	})
	if err != nil {
		panic(err)
	}

	path, err := discoverPathToMigrationsFolder("migrations")
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+path,
		config.Postgresql.Database,
		driver,
	)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			panic(err)
		}
	}

	return nil
}

func discoverPathToMigrationsFolder(folderName string) (string, error) {
	currentPath := folderName
	for tries := 10; tries > 0; tries-- {
		if _, err := os.Stat(currentPath); err == nil {
			currentPath = strings.ReplaceAll(currentPath, "\\", "//")
			return currentPath, nil
		}
		currentPath = filepath.Join("..", currentPath)
	}
	return "", fmt.Errorf("could not discover config")
}
