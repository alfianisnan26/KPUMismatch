package dao

import "kawalrealcount/internal/data/model"

type KPU interface {
	GetHHWCInfo(req *model.HHCWEntity) error
	GetPageLink(req []string) (string, error)
	GetPPWTList(entity model.PPWTEntity) ([]model.PPWTEntity, error)
}
