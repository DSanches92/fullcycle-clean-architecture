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
	migrationsPath       = "file://sql/migrations"
	maxRetries           = 30
	retryInterval        = 2 * time.Second
	dbPingRetries        = 15
	dbPingRetryInterval  = time.Second
)

// WaitForDB pings the database repeatedly until it accepts connections.
func WaitForDB(db *sql.DB) error {
	var lastErr error
	for i := 0; i < dbPingRetries; i++ {
		if err := db.Ping(); err == nil {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(dbPingRetryInterval)
	}
	return fmt.Errorf("database did not become ready after %d attempts: %w", dbPingRetries, lastErr)
}

// RunMigrations applies pending migrations with retry to handle race conditions
// between the database container startup and the application boot. It opens its
// own connection so the caller's *sql.DB is not closed by the migrate driver.
func RunMigrations(dsn string) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := runMigrationOnce(dsn); err == nil {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(retryInterval)
	}
	return fmt.Errorf("failed to apply migrations after %d attempts: %w", maxRetries, lastErr)
}

func runMigrationOnce(dsn string) error {
	db, err := sql.Open("mysql", dsn)
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

	m, err := migrate.NewWithInstance("file", srcDriver, "mysql", driver)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
