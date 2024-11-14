package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq" // Import the pq driver
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (db *Database) Close() error {
	return db.db.Close()
}
