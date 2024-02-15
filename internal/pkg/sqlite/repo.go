package sqlite

import (
	"context"
	"database/sql"
	"time"

	"kawalrealcount/internal/data/dao"

	_ "github.com/mattn/go-sqlite3"
)

type repo struct {
	db *sql.DB
}

type Param struct {
	FilePath string
}

func New(param Param) (dao.Cache, error) {
	db, err := sql.Open("sqlite3", param.FilePath)
	if err != nil {
		return nil, err
	}

	// Create the cache table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cache (
		key TEXT PRIMARY KEY,
		value BLOB
	)`)
	if err != nil {
		return nil, err
	}

	return &repo{
		db: db,
	}, nil
}

func (repo *repo) Get(key string, expiry time.Duration, fallback func(string) ([]byte, error)) ([]byte, error) {
	ctx := context.Background()

	// Try to get the value from SQLite cache
	var value []byte
	err := repo.db.QueryRowContext(ctx, "SELECT value FROM cache WHERE key = ?", key).Scan(&value)
	if err == nil {
		// Value found in cache, return it
		return value, nil
	} else if err != sql.ErrNoRows {
		// Some other error occurred
		return nil, err
	}

	// Key not found in cache, fetch from the fallback function
	fallbackVal, err := fallback(key)
	if err != nil {
		return nil, err
	}

	// Store the value in cache
	_, err = repo.db.ExecContext(ctx, "INSERT INTO cache(key, value) VALUES (?, ?)", key, fallbackVal)
	if err != nil {
		return nil, err
	}

	return fallbackVal, nil
}
