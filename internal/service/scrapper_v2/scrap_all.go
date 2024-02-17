package scrapper_v2

import (
	"fmt"
	"kawalrealcount/internal/data/model"
)

type hhcwCached struct {
	obj    *model.HHCWEntity
	cached bool
}

func (svc *service) ScrapAll() error {
	fmt.Println("Gathering PPWT")

	res, err := svc.getPpwtList()
	if err != nil {
		return err
	}

	hhcwCh := make(chan *hhcwCached)
	defer close(hhcwCh)

	fmt.Println("Gathering HHCW")
	go func() {
		svc.fetchHhcw(res, hhcwCh)
		hhcwCh <- nil
	}()

	webStats := model.NewWebStats(len(res))
	stats := model.Stats{
		Contributor: svc.Contributor.Email,
		WebStast:    webStats,
	}

	return svc.updateAll(hhcwCh, &stats)
}
