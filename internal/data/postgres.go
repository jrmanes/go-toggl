package data

import (
	"database/sql"
	"os"
)

type PostgresDB struct {
	db *sql.DB
}

func getConnection() (*sql.DB, error) {
	uri := os.Getenv("DATABASE_URI")
	return sql.Open("postgres", uri)
}
