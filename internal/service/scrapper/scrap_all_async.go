package scrapper

import (
	"fmt"
	"kawalrealcount/internal/data/model"
)

func (svc *service) ScrapAllAsyncSeedOnlyByMap(ppwtMap map[string]*model.PPWTEntity) error {

	hhcwCh := make(chan *model.HHCWEntity, len(ppwtMap))

	var (
		webStats = model.NewWebStats(len(ppwtMap))
		stats    = model.Stats{
			Contributor: svc.contributor.Email,
			WebStast:    webStats,
		}
	)

	if err := svc.databaseRepo.InsertStats(&stats); err != nil {
		return err
	}

	var finished bool
	var duplicateCount int

	go svc.fetchHhcwByMap(ppwtMap, hhcwCh)

	groupLength := svc.batchInsertLength
	group := make(map[string]*model.HHCWEntity, groupLength)

	for !finished {
		hhcw := <-hhcwCh
		if hhcw == nil {
			finished = true
		} else if _, found := group[hhcw.Parent.Kode]; !found {
			group[hhcw.Parent.Kode] = hhcw
			stats.Evaluate(hhcw)
		} else {
			duplicateCount++
		}

		if len(group) >= groupLength || (finished && len(group) > 0) {
			if err := webStats.Update(len(group)); err != nil {
				fmt.Println(err)
			}

			var firstEntity *model.HHCWEntity
			for _, entity := range group {
				firstEntity = entity
				break
			}

			fmt.Printf("%v\t| Dupl: %d | Sample: %s | %s\n", webStats.String(), duplicateCount, firstEntity.Parent.Kode, firstEntity.Link)

			if err := svc.databaseRepo.PutReplaceMultipleData(group, stats.WebStast.UploadID); err != nil {
				fmt.Println(err.Error())
			}

			if err := svc.databaseRepo.UpdateStats(&stats); err != nil {
				fmt.Println(err.Error())
			}

			if err := svc.databaseRepo.InsertWebStats(webStats); err != nil {
				fmt.Println(err.Error())
			}

			for key := range group {
				delete(group, key)
			}
		}
	}

	stats.Finalize()

	fmt.Printf("%v\t| HHCW FINISHED\n", stats.WebStast.String())
	return svc.databaseRepo.UpdateStats(&stats)
}
