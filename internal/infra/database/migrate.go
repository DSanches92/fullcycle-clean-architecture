package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	filesource "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/go-sql-driver/mysql"
)

const (
	migrationsPath      = "file://sql/migrations"
	maxRetries          = 30
	retryInterval       = 2 * time.Second
	dbPingRetries       = 15
	dbPingRetryInterval = time.Second
)

func WaitForDB(db *sql.DB) error {
	var lastErr error

	for range dbPingRetries {
		if err := db.Ping(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		time.Sleep(dbPingRetryInterval)
	}

	return fmt.Errorf("banco de dados não ficou pronto após %d tentativas: %w", dbPingRetries, lastErr)
}

func RunMigrations(dbDriver, connString string) error {
	var lastErr error

	for range maxRetries {
		if err := runMigrationOnce(dbDriver, connString); err == nil {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(retryInterval)
	}

	return fmt.Errorf("falha ao aplicar as migrations após %d tentativas: %w", maxRetries, lastErr)
}

func runMigrationOnce(dbDriver, connString string) error {
	db, err := sql.Open(dbDriver, connString)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := migratemysql.WithInstance(db, &migratemysql.Config{})
	if err != nil {
		return err
	}

	var fileSource filesource.File
	srcDriver, err := fileSource.Open(migrationsPath)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("file", srcDriver, dbDriver, driver)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
