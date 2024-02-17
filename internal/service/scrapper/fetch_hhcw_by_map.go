package scrapper

import (
	"fmt"
	"kawalrealcount/internal/data/model"
	"kawalrealcount/internal/pkg/semaphore"
	"sync"
)

func (svc *service) fetchHhcwByMap(ppwtMap map[string]*model.PPWTEntity, hhcwCh chan<- *model.HHCWEntity) {
	fmt.Println("Start To Fetch HHCW:", len(ppwtMap))
	sm := semaphore.NewSemaphore(svc.maximumRunningThread)
	var wg sync.WaitGroup
	wg.Add(len(ppwtMap))

	for _, ppwt := range ppwtMap {
		sm.Acquire()
		go func(ppwt *model.PPWTEntity, hhcwCh chan<- *model.HHCWEntity) {
			defer sm.Release()
			defer wg.Done()

			res, err := svc.kpuRepo.GetHHCWInfo(ppwt)
			if err != nil {
				fmt.Println(err)
				return
			}

			res.Link, _ = svc.kpuRepo.GetPageLink(*res.Parent)
			hhcwCh <- &res
		}(ppwt, hhcwCh)
	}

	wg.Wait()
	hhcwCh <- nil
	return
}
