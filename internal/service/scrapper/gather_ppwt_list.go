package scrapper

import (
	"fmt"
	"kawalrealcount/internal/data/model"
	"time"
)

func (svc *service) GatherAllPPWTMap() (map[string]*model.PPWTEntity, error) {
	ppwtCh := make(chan *model.PPWTEntity)

	go func() {
		startingPoint := model.NewPPWT("0")
		if err := svc.recursiveFetchPpwt(&startingPoint, ppwtCh); err != nil {
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
