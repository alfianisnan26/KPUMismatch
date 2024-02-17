package scrapper_v2

import (
	"kawalrealcount/internal/data/dao"
	"kawalrealcount/internal/data/model"
	"time"
)

type Service interface {
	ScrapAll() error
}

type Param struct {
	KPURepo                    dao.KPU
	CacheRepo                  dao.Cache
	DatabaseRepo               dao.Database
	MaximumRunningThread       int
	BatchInsertLength          int
	Contributor                model.ContributorData
	ValidRecordExpiry          time.Duration
	NotNullInvalidRecordExpiry time.Duration
	NullRecordExpiry           time.Duration
}

type service struct {
	Param
}

func New(param Param) Service {
	return &service{
		Param: param,
	}
}
