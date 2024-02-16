package postgresql

import (
	"fmt"
	"kawalrealcount/internal/data/model"
	"time"
)

func (r *repo) InsertStats(stats *model.Stats) error {
	stats.CreatedAt = time.Now()

	query := fmt.Sprintf(`INSERT INTO %s (created_at) VALUES ($1) RETURNING id;`, r.tableStat)
	if err := r.db.QueryRow(query,
		stats.CreatedAt.UTC().UnixMilli(),
	).Scan(&(stats.ID)); err != nil {
		return err
	}

	fmt.Println("New record ID is:", stats.ID)
	return nil
}