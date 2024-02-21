package scrapper_v2

import (
	"fmt"
	"kawalrealcount/internal/data/model"
	"sync"
)

const targetLength = 823236

type hhcwCached struct {
	obj       *model.HHCWEntity
	cached    bool
	isChanged bool
}

func (svc *service) ScrapAll() error {
	fmt.Println("Gathering PPWT")

	ppwtCh := make(chan *model.HHCWEntity, svc.MaximumRunningThread)
	defer close(ppwtCh)
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	go func() {
		wg.Done()
		err := svc.DatabaseRepo.GetPPWTCodeOnly(ppwtCh)
		if err != nil {
			return
		}

		ppwtCh <- nil
	}()

	hhcwCh := make(chan *hhcwCached, svc.MaximumRunningThread)
	defer close(hhcwCh)

	fmt.Println("Gathering HHCW")
	wg.Add(1)
	go func() {
		defer wg.Done()
		svc.fetchHhcw(ppwtCh, hhcwCh)
		hhcwCh <- nil
	}()

	webStats := model.NewWebStats(targetLength)
	stats := model.Stats{
		Contributor: svc.Contributor.Email,
		WebStast:    webStats,
	}

	return svc.updateAll(hhcwCh, &stats)
}
