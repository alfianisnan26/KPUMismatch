package scrapper

import (
	"kawalrealcount/internal/data/dao"
	"time"
)

type Service interface {
	ScrapAllCompiled(filePath string, isScrapAll bool) error
	ScrapAllSeedOnly() error
}

type Param struct {
	KPURepo              dao.KPU
	KawalPemiluRepo      dao.KawalPemilu
	DatabaseRepo         dao.Database
	MaximumRunningThread int
	ProgressRefreshRate  time.Duration
}

type service struct {
	kpuRepo              dao.KPU
	kawalPemiluRepo      dao.KawalPemilu
	databaseRepo         dao.Database
	maximumRunningThread int
	progressRefreshRate  time.Duration
}

func New(param Param) Service {
	return &service{
		kpuRepo:              param.KPURepo,
		kawalPemiluRepo:      param.KawalPemiluRepo,
		maximumRunningThread: param.MaximumRunningThread,
		databaseRepo:         param.DatabaseRepo,
		progressRefreshRate:  param.ProgressRefreshRate,
	}
}
