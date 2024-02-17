package postgresql

import (
	"database/sql"
	_ "github.com/lib/pq"
	"kawalrealcount/internal/data/dao"
)

type repo struct {
	db                       *sql.DB
	tableRecord              string
	tableStat, tableWebStats string
}

type Param struct {
	ConnectionURL string
	TableRecord   string
	TableStats    string
	TableWebStats string
}

func New(param Param) (dao.Database, error) {
	// Initialize PostgreSQL connection
	db, err := sql.Open("postgres", param.ConnectionURL)
	if err != nil {
		return nil, err
	}

	// Check if the connection is successful
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Return the PostgreSQL repository
	return &repo{
		db:            db,
		tableRecord:   param.TableRecord,
		tableStat:     param.TableStats,
		tableWebStats: param.TableWebStats,
	}, nil
}
