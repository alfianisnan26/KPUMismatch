package updater

import "fmt"

func (s service) UpdateStats() error {
	summary, err := s.UpdaterDatabaseRepo.GetSummary()
	if err != nil {
		return err
	}

	fmt.Println("Updating Status...")

	return s.UpdaterDatabaseRepo.InsertSummary(summary)
}
