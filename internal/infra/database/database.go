package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/DSanches92/fullcycle-clean-architecture/configs"
)

func DatabaseConn(env configs.Database) (*sql.DB, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		env.DBUser,
		env.DBPassword,
		env.DBHost,
		env.DBPort,
		env.DBName,
	)

	connection, err := sql.Open(env.DBDriver, connString)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir a conexão com banco de dados: %v", err)
	}

	connection.SetMaxOpenConns(15)
	connection.SetMaxIdleConns(3)
	connection.SetConnMaxLifetime(2 * time.Minute)
	connection.SetConnMaxIdleTime(30 * time.Second)

	if err := WaitForDB(connection); err != nil {
		return nil, fmt.Errorf("banco de dados não está pronto: %v", err)
	}

	if err := RunMigrations(env.DBDriver, connString); err != nil {
		return nil, fmt.Errorf("falha ao rodar as migrations: %v", err)
	}

	return connection, nil
}
