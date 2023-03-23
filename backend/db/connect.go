package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	// connStr := "postgresql://<username>:<password>@<database_ip>/todos?sslmode=disable"
	dsn := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_DB_USERNAME"),
		os.Getenv("POSTGRES_DB_PASSWORD"),
		os.Getenv("POSTGRES_DB_HOST"),
		os.Getenv("POSTGRES_DB_NAME"),
		os.Getenv("POSTGRES_DB_SSL_MODE"))
	db, errOpen := sql.Open("postgres", dsn)
	if errOpen != nil {
		return nil, errOpen
	}

	return db, nil
}
