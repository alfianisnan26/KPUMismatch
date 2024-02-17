package scrapper

import (
	"fmt"
	"kawalrealcount/internal/data/model"
)

func (svc *service) ScrapAllSyncSeedOnlyByMap(ppwtMap map[string]*model.PPWTEntity) error {
	group := make(map[string]*model.HHCWEntity, svc.batchInsertLength)
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

	for k, v := range ppwtMap {
		if len(group) >= svc.batchInsertLength {
			if err := svc.upload(&group, &stats); err != nil {
				return err
			}

			for key := range group {
				delete(group, key)
			}
		}

		res, err := svc.kpuRepo.GetHHCWInfo(v)
		if err != nil {
			continue
		}

		if k != res.Parent.Kode {
			resLink, _ := svc.kpuRepo.GetPageLink(*res.Parent)
			refLink, _ := svc.kpuRepo.GetPageLink(*v)

			fmt.Printf("UNMATCHED | Ref: [%s|%s] %s | Res: [%s] %s\n", k, v.Kode, refLink, res.Parent.Kode, resLink)
			return nil
		}

		group[k] = &res
	}

	stats.Finalize()

	if err := svc.upload(&group, &stats); err != nil {
		return err
	}

	return nil
}

func (svc *service) upload(group *map[string]*model.HHCWEntity, stats *model.Stats) error {
	if err := stats.WebStast.Update(len(*group)); err != nil {
		return err
	}

	var firstEntity *model.HHCWEntity
	for _, entity := range *group {
		firstEntity = entity
		break
	}

	fmt.Printf("%v\t | Sample: %s | %s\n", stats.WebStast.String(), firstEntity.Parent.Kode, firstEntity.Link)

	err := svc.databaseRepo.PutReplaceMultipleData(*group, stats.WebStast.UploadID)
	if err != nil {
		return err
	}

	if err := svc.databaseRepo.UpdateStats(stats); err != nil {
		return err
	}

	return svc.databaseRepo.InsertWebStats(stats.WebStast)

}
