package scrapper

import (
	"kawalrealcount/internal/data/model"
)

func (svc *service) ScrapPPWTWithStartingPoint(entity model.PPWTEntity) ([]model.PPWTEntity, error) {
	if entity.LowestLevel() {
		return []model.PPWTEntity{entity}, nil
	}

	res, err := svc.kpuRepo.GetPPWTList(entity)
	if err != nil {
		return nil, err
	}

	out := make([]model.PPWTEntity, 0)
	for _, re := range res {
		res, err := svc.ScrapPPWTWithStartingPoint(re)
		if err != nil {
			return nil, err
		}

		out = append(out, res...)
	}

	return out, nil
}
