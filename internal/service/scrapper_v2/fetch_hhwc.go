package scrapper_v2

import (
	"kawalrealcount/internal/data/model"
	"kawalrealcount/internal/pkg/semaphore"
	"time"
)

func (svc *service) fetchHhcw(ppwtCh <-chan *model.HHCWEntity, hhcwCh chan<- *hhcwCached) {
	sm := semaphore.NewSemaphore(svc.MaximumRunningThread)

	for {
		entity := <-ppwtCh
		if entity == nil {
			break
		}
		sm.Acquire()
		go func(entity *model.HHCWEntity) {
			defer sm.Release()

			switch {
			case entity.IsNonNullVote() && time.Since(entity.ObtainedAt) < svc.NotNullInvalidRecordExpiry:
				fallthrough
			case entity.IsValidVote() && time.Since(entity.ObtainedAt) < svc.ValidRecordExpiry:
				fallthrough
			case time.Since(entity.ObtainedAt) < svc.NullRecordExpiry:
				hhcwCh <- &hhcwCached{
					obj:    entity,
					cached: true,
				}
				return
			}

			old := *entity

			if err := svc.KPURepo.GetHHWCInfo(entity); err != nil {
				return
			}

			hhcwCh <- &hhcwCached{
				obj:       entity,
				isChanged: entity.UpdatedAt != old.UpdatedAt,
			}
		}(entity)
	}

	for i := 0; i < svc.MaximumRunningThread; i++ {
		sm.Acquire()
	}
}
