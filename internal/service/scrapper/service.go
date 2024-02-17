package scrapper

import (
	"kawalrealcount/internal/data/dao"
	"kawalrealcount/internal/data/model"
)

type Service interface {
	ScrapAllCompiled(filePath string, isScrapAll bool) error
	ScrapAllAsyncSeedOnlyByMap(map[string]*model.PPWTEntity) error
	ScrapAllSyncSeedOnlyByMap(map[string]*model.PPWTEntity) error
	GatherAllPPWTMap() (map[string]*model.PPWTEntity, error)
}

type Param struct {
	KPURepo              dao.KPU
	KawalPemiluRepo      dao.KawalPemilu
	DatabaseRepo         dao.Database
	MaximumRunningThread int
	Contributor          model.ContributorData
	BatchInsertLength    int
}

type service struct {
	kpuRepo              dao.KPU
	kawalPemiluRepo      dao.KawalPemilu
	databaseRepo         dao.Database
	maximumRunningThread int
	contributor          model.ContributorData
	batchInsertLength    int
}

func New(param Param) *service {
	return &service{
		kpuRepo:              param.KPURepo,
		kawalPemiluRepo:      param.KawalPemiluRepo,
		maximumRunningThread: param.MaximumRunningThread,
		databaseRepo:         param.DatabaseRepo,
		contributor:          param.Contributor,
		batchInsertLength:    param.BatchInsertLength,
	}
}
