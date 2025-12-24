package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kiart-tantasi/crm-go/internal/env"
)

func Connect() (*sql.DB, error) {
	dbUser, err := env.GetEnvRequired("DB_USER")
	if err != nil {
		return nil, err
	}
	dbPass, err := env.GetEnvRequired("DB_PASS")
	if err != nil {
		return nil, err
	}
	dbHost, err := env.GetEnvRequired("DB_HOST")
	if err != nil {
		return nil, err
	}
	dbPort, err := env.GetEnvRequired("DB_PORT")
	if err != nil {
		return nil, err
	}
	dbName, err := env.GetEnvRequired("DB_NAME")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// db.Ping helps verify connection is valid (e.g. credentials are correct)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")
	return db, nil
}
