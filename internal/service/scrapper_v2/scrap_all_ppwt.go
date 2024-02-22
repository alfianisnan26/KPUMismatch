package scrapper_v2

import (
	"fmt"
	"sync"

	"kawalrealcount/internal/data/model"
)

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

func (svc *service) ScrapAllPPWT() error {
	var ppwtCh = make(chan *model.PPWTEntity)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := svc.recursivePpwt(model.NewPPWT("0"), ppwtCh); err != nil {
			fmt.Println(err.Error())
		}

		ppwtCh <- nil
	}()

	mapPpwt := make(map[string]*model.PPWTEntity)
	for finish := false; !finish; {
		ppwt := <-ppwtCh
		if ppwt == nil {
			finish = true
		} else if v, found := mapPpwt[ppwt.Kode]; found {
			fmt.Println("Already Found:", v.Kode, ppwt.Kode)
		} else {
			mapPpwt[ppwt.Kode] = ppwt
		}

		if (len(mapPpwt) > 0 && len(mapPpwt)%svc.BatchInsertLength == 0) || (finish && len(mapPpwt) > 0) {

			fmt.Printf("[%d] PPWT Gathered Sample: %s\n", len(mapPpwt), ppwt.GetCanonicalCode())

			if err := svc.DatabaseRepo.PutReplacePPWT(mapPpwt); err != nil {
				return err
			}

			for k := range mapPpwt {
				delete(mapPpwt, k)
			}
		}
	}

	wg.Wait()

	return nil
}
