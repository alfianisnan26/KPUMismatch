package scrapper

import (
	"fmt"
	"kawalrealcount/internal/data/model"
	"kawalrealcount/internal/pkg/semaphore"
	"time"
)

const targetTotalTps = 823236

func (svc *service) ScrapAllSeedOnly() error {

	var (
		stats model.Stats
	)

	if err := svc.databaseRepo.InsertStats(&stats); err != nil {
		return err
	}

	ppwtCh := make(chan model.PPWTEntity)
	errCh := make(chan error)

	go func() {
		if err := svc.recursiveFetchPpwt(model.NewPPWT("0"), ppwtCh); err != nil {
			errCh <- err
		}

		ppwtCh <- model.PPWTEntity{}
	}()
	var start = time.Now()
	var nextLogTime time.Time
	var count int
	sm := semaphore.NewSemaphore(svc.maximumRunningThread)

LoopFor:
	for {
		select {
		case ppwt := <-ppwtCh:
			if ppwt.ID == 0 {
				// eof
				break LoopFor
			}

			link, _ := svc.kpuRepo.GetPageLink(ppwt)

			if tn := time.Now(); tn.After(nextLogTime) {
				nextLogTime = tn.Add(svc.progressRefreshRate)

				percentage := float32(count) / float32(targetTotalTps) * 100
				est := estimateTime(start, percentage, 100)
				stats.Update(percentage, est, start)

				if err := svc.databaseRepo.UpdateStats(&stats); err != nil {
					fmt.Println(err)
				}

				fmt.Printf("[%2.2f%%] %s  Est: (minute) %2.2f\t%s\n", percentage, ppwt.Kode, est.Minutes(), link)

			}

			sm.Acquire()
			go func(ppwt model.PPWTEntity, link string, count int, stats *model.Stats) {
				defer sm.Release()

				err := svc.processPpwt(ppwt, link, stats)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

			}(ppwt, link, count, &stats)

			count++
		}
	}

	stats.Finalize(start)
	return svc.databaseRepo.UpdateStats(&stats)
}

func (svc *service) processPpwt(entity model.PPWTEntity, link string, stats *model.Stats) error {
	res, err := svc.kpuRepo.GetHHCWInfo(entity)
	if err != nil {
		return err
	}

	if res.Administrasi.Suara.Total == 0 || res.Chart.Sum() == 0 {
		// skipped
		return nil
	}

	res.Link = link

	stats.Evaluate(res)
	if err := svc.databaseRepo.PutReplaceData(res); err != nil {
		return err
	}

	return nil
}

func (svc *service) recursiveFetchPpwt(entity model.PPWTEntity, ppwtCh chan<- model.PPWTEntity) error {
	if entity.LowestLevel() {
		ppwtCh <- entity
		return nil
	}

	res, err := svc.kpuRepo.GetPPWTList(entity)
	if err != nil {
		return err
	}

	for _, re := range res {
		if err := svc.recursiveFetchPpwt(re, ppwtCh); err != nil {
			return err
		}
	}

	return nil
}

// estimateTime calculates the estimated time until reaching the target percentage.
func estimateTime(startTime time.Time, currentPercentage, targetPercentage float32) time.Duration {

	elapsed := time.Since(startTime)
	remainingPercentage := targetPercentage - currentPercentage
	ct := time.Duration(currentPercentage * 1000)
	rt := time.Duration(remainingPercentage * 1000)
	if ct == 0 || rt == 0 {
		return 0
	}

	remainingTime := elapsed / ct * rt
	return remainingTime
}
