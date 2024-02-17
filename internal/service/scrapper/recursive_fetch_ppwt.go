package scrapper

import "kawalrealcount/internal/data/model"

func (svc *service) recursiveFetchPpwt(entity *model.PPWTEntity, ppwtCh chan<- *model.PPWTEntity) error {
	if entity.LowestLevel() {
		ppwtCh <- entity
		return nil
	}

	res, err := svc.kpuRepo.GetPPWTList(*entity)
	if err != nil {
		return err
	}

	for _, re := range res {
		if err := svc.recursiveFetchPpwt(&re, ppwtCh); err != nil {
			return err
		}
	}

	return nil
}
