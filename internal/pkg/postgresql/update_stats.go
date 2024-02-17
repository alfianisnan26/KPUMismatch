package postgresql

import (
	"fmt"
	"kawalrealcount/internal/data/model"
)

func (r *repo) UpdateStats(stats *model.Stats) error {
	query := `UPDATE %s
				SET
					progress = $1,
					estimate_time = $2,
					last_progress_update = $3,
					processing_time = $4,
					finished_at = $5,
				    total_record = $6
				WHERE id = $7;`

	var finishedAt int64
	if !stats.FinishedAt.IsZero() {
		finishedAt = stats.FinishedAt.UTC().UnixMilli()
	}

	_, err := r.db.Exec(fmt.Sprintf(query, r.tableStat),
		stats.WebStast.Percentage,                    //1
		stats.WebStast.Estimation.Milliseconds(),     //2
		stats.WebStast.Timestamp.UTC().UnixMilli(),   //3
		stats.WebStast.ProcessingTime.Milliseconds(), //4
		finishedAt, //5
		stats.WebStast.DataCount,
		stats.WebStast.UploadID,
	)

	return err
}
