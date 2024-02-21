package updater

import (
	"encoding/json"
	"fmt"
	"kawalrealcount/internal/data/model"
	"time"
)

func (s service) UpdateStats() error {
	summary, err := s.UpdaterDatabaseRepo.GetSummary()
	if err != nil {
		return err
	}

	fmt.Printf("[%v] Updating Status...\n", time.Now().String())

	return s.UpdaterDatabaseRepo.InsertSummary(summary)
}

func (s service) UpdateStaticStats() error {
	mapDist, err := s.UpdaterDatabaseRepo.GetMapDist()
	if err != nil {
		return err
	}

	mapDistBuf, err := json.Marshal(mapDist)
	if err != nil {
		return err
	}

	tab1, err := s.UpdaterDatabaseRepo.GetPotentialTableSum()
	if err != nil {
		return err
	}

	tab1Buf, err := json.Marshal(tab1)
	if err != nil {
		return err
	}

	tab2, err := s.UpdaterDatabaseRepo.GetPotentialTableAllIn()
	if err != nil {
		return err
	}

	tab2Buf, err := json.Marshal(tab2)
	if err != nil {
		return err
	}

	summaries := []model.StaticSummary{
		{
			Key: "map_dist",
			Val: mapDistBuf,
		},
		{
			Key: "potential_table_sum",
			Val: tab1Buf,
		},
		{
			Key: "potential_table_all_in",
			Val: tab2Buf,
		},
	}

	fmt.Printf("[%v] Updating Static Summary...\n", time.Now().String())

	return s.UpdaterDatabaseRepo.UpdateStaticSummary(summaries)
}
