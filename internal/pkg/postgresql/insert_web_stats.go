package postgresql

import (
	"encoding/json"
	"fmt"
	"kawalrealcount/internal/data/model"
)

func (r *repo) InsertWebStats(webStats *model.WebStats) error {
	memory, err := json.Marshal(webStats.MemoryInfo)
	if err != nil {
		return err
	}

	query := `INSERT INTO %s (upload_id, memory_info, rps, est, percentage, data_count, ts)
values ($1,$2,$3,$4,$5,$6,$7);`

	_, err = r.db.Exec(fmt.Sprintf(query, r.tableWebStats),
		webStats.UploadID,
		memory,
		webStats.RPS,
		webStats.Estimation.Milliseconds(),
		webStats.Percentage,
		webStats.DataCount,
		webStats.Timestamp.UTC().UnixMilli(),
	)
	return err
}
