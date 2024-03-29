package scrapper_v2

import (
	"fmt"
	"kawalrealcount/internal/data/model"
)

func (svc *service) updateAll(hhcwCh <-chan *hhcwCached, stats *model.Stats) error {

	err := svc.DatabaseRepo.InsertStats(stats)
	if err != nil {
		return err
	}

	fmt.Println("Updating Stats", stats.WebStast.UploadID)

	group := make([]*model.HHCWEntity, 0, svc.BatchInsertLength)
	var finished bool
	var cachedCount, changedCount int

	for !finished {
		hhcw := <-hhcwCh
		if hhcw != nil {
			if hhcw.cached {
				cachedCount++
			} else {
				if hhcw.obj == nil {
					fmt.Println("Unknown Nil Object", hhcw)
					continue
				}

				if hhcw.isChanged {
					changedCount++
				}

				hhcw.obj.Link, _ = svc.KPURepo.GetPageLink(hhcw.obj.GetCanonicalCode())
				group = append(group, hhcw.obj)
			}

			stats.Evaluate()
		} else {
			finished = true
		}

		if len(group) >= svc.BatchInsertLength || finished && len(group) > 0 {
			if err := stats.WebStast.Update(len(group)); err != nil {
				fmt.Println(err.Error())
			}

			var firstEntity *model.HHCWEntity
			for _, entity := range group {
				firstEntity = entity
				break
			}

			fmt.Printf("%v\t| Changed: %d | Cached: %d | Sample: %s | %s\n", stats.WebStast.String(), changedCount, cachedCount, firstEntity.Code, firstEntity.Link)

			if err := svc.DatabaseRepo.PutReplaceListData(group, stats.WebStast.UploadID); err != nil {
				fmt.Println(err.Error())
			}

			if err := svc.DatabaseRepo.UpdateStats(stats); err != nil {
				fmt.Println(err.Error())
			}

			if err := svc.DatabaseRepo.InsertWebStats(stats.WebStast); err != nil {
				fmt.Println(err.Error())
			}

			group = group[:0]
		}

	}

	if finished {
		stats.Finalize()

		fmt.Printf("%v\t| HHCW FINISHED\n", stats.WebStast.String())
		return svc.DatabaseRepo.UpdateStats(stats)
	}

	return nil
}
