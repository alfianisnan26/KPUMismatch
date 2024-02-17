package scrapper_v2

import (
	"fmt"
	"kawalrealcount/internal/data/model"
	"time"
)

func (svc *service) getPpwtList() (map[string]*model.PPWTEntity, error) {
	var ppwtCh = make(chan *model.PPWTEntity)

	go func() {
		if err := svc.recursivePpwt(model.NewPPWT("0"), ppwtCh); err != nil {
			fmt.Println(err.Error())
		}

		ppwtCh <- nil
	}()

	mapPpwt := make(map[string]*model.PPWTEntity)

	for {
		ppwt := <-ppwtCh
		if ppwt == nil {
			break
		}

		if len(mapPpwt)%10000 == 0 {
			fmt.Printf("[%d] PPWT Gathered Sample: %s\n", len(mapPpwt), ppwt.GetCanonicalCode())
		}

		if v, found := mapPpwt[ppwt.Kode]; found {
			fmt.Println("Already Found:", v.Kode, ppwt.Kode)
			time.Sleep(time.Second * 10)
		}
		mapPpwt[ppwt.Kode] = ppwt
	}

	return mapPpwt, nil
}

func (svc *service) recursivePpwt(ppwt model.PPWTEntity, ppwtCh chan<- *model.PPWTEntity) error {
	if ppwt.LowestLevel() {
		ppwtCh <- &ppwt
		return nil
	}

	res, err := svc.KPURepo.GetPPWTList(ppwt)
	if err != nil {
		return err
	}

	for _, entity := range res {
		if err := svc.recursivePpwt(entity, ppwtCh); err != nil {
			fmt.Println(err.Error())
			// ignore error here
		}
	}

	return nil
}
