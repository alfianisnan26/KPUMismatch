package scrapper

import (
	"kawalrealcount/internal/data/dao"
)

type Service interface {
	ScrapAllCompiled(filePath string) error
}

type Param struct {
	KPURepo              dao.KPU
	KawalPemiluRepo      dao.KawalPemilu
	DatabaseRepo         dao.Database
	MaximumRunningThread int
}

type service struct {
	kpuRepo              dao.KPU
	kawalPemiluRepo      dao.KawalPemilu
	databaseRepo         dao.Database
	maximumRunningThread int
}

func New(param Param) Service {
	return &service{
		kpuRepo:              param.KPURepo,
		kawalPemiluRepo:      param.KawalPemiluRepo,
		maximumRunningThread: param.MaximumRunningThread,
		databaseRepo:         param.DatabaseRepo,
	}
}
