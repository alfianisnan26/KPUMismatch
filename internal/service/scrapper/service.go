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
	MaximumRunningThread int
}

type service struct {
	kpuRepo              dao.KPU
	kawalPemiluRepo      dao.KawalPemilu
	maximumRunningThread int
}

func New(param Param) Service {
	return &service{
		kpuRepo:              param.KPURepo,
		kawalPemiluRepo:      param.KawalPemiluRepo,
		maximumRunningThread: param.MaximumRunningThread,
	}
}
