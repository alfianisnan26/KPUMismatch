package scrapper_v2

import (
	"time"

	"kawalrealcount/internal/data/dao"
	"kawalrealcount/internal/data/model"
)

type Service interface {
	ScrapAll() error
	ScrapAllPPWT() error
}

type Param struct {
	KPURepo                    dao.KPU
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
