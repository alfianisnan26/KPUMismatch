package updater

import "kawalrealcount/internal/data/dao"

type service struct {
	Param
}

type Service interface {
	UpdateStats() error
	UpdateStaticStats() error
}

type Param struct {
	UpdaterDatabaseRepo dao.UpdaterDatabase
}

func New(param Param) Service {
	return service{
		Param: param,
	}
}
