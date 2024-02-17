package scrapper_v2

import (
	"fmt"
	"kawalrealcount/internal/data/model"
	"kawalrealcount/internal/pkg/semaphore"
	"sync"
	"time"
)

func (svc *service) fetchHhcw(ppwtMap map[string]*model.PPWTEntity, hhcwCh chan<- *hhcwCached) {
	sm := semaphore.NewSemaphore(svc.MaximumRunningThread)
	var wg sync.WaitGroup
	wg.Add(len(ppwtMap))
	for _, entity := range ppwtMap {
		sm.Acquire()
		go func(entity *model.PPWTEntity) {
			defer wg.Done()
			defer sm.Release()

			// if already available in cache
			cached, err := svc.CacheRepo.GetHHCW(entity.Kode)
			if err == nil {
				hhcwCh <- &hhcwCached{
					obj:    &cached,
					cached: true,
				}
				return
			}

			res, err := svc.KPURepo.GetHHWCNoCacheInfo(entity)
			if err != nil {
				return
			}

			var expiry time.Duration
			if res.IsNonNullVote() {
				if res.IsValidVote() {
					expiry = svc.ValidRecordExpiry // 2 h
				}
				expiry = svc.NotNullInvalidRecordExpiry // 30 min
			} else {
				expiry = svc.NullRecordExpiry // 10 min
			}

			if err := svc.CacheRepo.PutHHCW(entity.Kode, res, expiry); err != nil {
				fmt.Println(err.Error())
			}

			hhcwCh <- &hhcwCached{
				obj: &res,
			}
		}(entity)
	}

	wg.Wait()
}
