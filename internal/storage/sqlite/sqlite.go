package sqlite

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"fmt"
)

type Storage struct {
	DB *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "internal.storage.sqlite.New"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	query := `
		CREATE TABLE IF NOT EXISTS Tasks (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Name TEXT NOT NULL UNIQUE,
			Completed INTEGER,
			Created_at DATETIME,
			Completed_at DATETIME
		) `

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db}, nil
}