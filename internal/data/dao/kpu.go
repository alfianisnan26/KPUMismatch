package dao

import "kawalrealcount/internal/data/model"

type KPU interface {
	GetPPWTList(req model.PPWTEntity) ([]model.PPWTEntity, error)
	GetHHCWInfo(req model.PPWTEntity) (model.HHCWEntity, error)
	GetPageLink(req model.PPWTEntity) (string, error)
	GetPPWTParent(req model.PPWTEntity) (model.PPWTEntity, error)
}
