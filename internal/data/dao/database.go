package dao

import "kawalrealcount/internal/data/model"

type Database interface {
	PutReplaceData(entity model.HHCWEntity, updateId uint64) error
	InsertStats(stats *model.Stats) error
	UpdateStats(stats *model.Stats) error
}
