package dao

import "kawalrealcount/internal/data/model"

type Database interface {
	PutReplaceListData(entities []*model.HHCWEntity, updateId uint64) error
	PutReplacePPWT(entities map[string]*model.PPWTEntity) error
	InsertStats(stats *model.Stats) error
	UpdateStats(stats *model.Stats) error
	InsertWebStats(webStats *model.WebStats) error
	GetPPWTCodeOnly(chan<- *model.HHCWEntity) error
}

type UpdaterDatabase interface {
	GetSummary() (model.Summary, error)
	UpdateStaticSummary(summaries []model.StaticSummary) error
	InsertSummary(summary model.Summary) error

	GetMapDist() ([]model.MapDist, error)
	GetPotentialTableSum() ([]model.HHCWEntity, error)
	GetPotentialTableAllIn() ([]model.HHCWEntity, error)
}
