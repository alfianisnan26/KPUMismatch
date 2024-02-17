package dao

import "kawalrealcount/internal/data/model"

type Database interface {
	PutReplaceMultipleData(entities map[string]*model.HHCWEntity, updateId uint64) error
	PutReplaceListData(entities []*model.HHCWEntity, updateId uint64) error
	InsertStats(stats *model.Stats) error
	UpdateStats(stats *model.Stats) error
	InsertWebStats(webStats *model.WebStats) error
}
