package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("POSTGRES_HOST", "localhost"),
		getEnv("POSTGRES_PORT", "5432"),
		getEnv("POSTGRES_USER", "postgres"),
		getEnv("POSTGRES_PASSWORD", ""),
		getEnv("POSTGRES_DB", "postgres"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("[db] error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("[db] error connecting to database: %w, check status or credential", err)
	}

	return db, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
