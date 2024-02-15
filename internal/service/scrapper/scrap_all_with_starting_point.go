package scrapper

import (
	"kawalrealcount/internal/data/model"
	"kawalrealcount/internal/pkg/semaphore"
	"sync"
)

func (svc *service) ScrapAllWithStartingPoint(hhcwListCh chan<- model.HHCWEntity, startingPoint model.PPWTEntity, ppwtList []model.PPWTEntity) error {
	var err error
	if startingPoint.Tingkat != 0 && startingPoint.Parent == nil {
		startingPoint, err = svc.kpuRepo.GetPPWTParent(startingPoint)
		if err != nil {
			return err
		}
	}

	go func() {
		wg := new(sync.WaitGroup)
		sm := semaphore.NewSemaphore(svc.maximumRunningThread)

		wg.Add(len(ppwtList))

		for _, entity := range ppwtList {

			sm.Acquire()
			go func(entity model.PPWTEntity) {
				defer sm.Release()
				defer wg.Done()

				res, err := svc.kpuRepo.GetHHCWInfo(entity)
				if err != nil {
					return
				}

				hhcwListCh <- res
			}(entity)
		}

		wg.Wait()
		close(hhcwListCh)
	}()

	return nil
}
