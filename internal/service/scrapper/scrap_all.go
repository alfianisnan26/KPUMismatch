package scrapper

import (
	"encoding/csv"
	"fmt"
	"kawalrealcount/internal/data/model"
	"kawalrealcount/internal/pkg/semaphore"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func (svc *service) ScrapAll(criterion model.Criterion, filePath string) error {
	startingPoint := model.NewPPWT("0")

	res, err := svc.kpuRepo.GetPPWTList(startingPoint)
	if err != nil {
		return err
	}

	for i, re := range res {
		fileName := filePath
		ext := filepath.Ext(fileName)
		fileName = strings.TrimSuffix(fileName, ext)
		fileName = fileName + "_" + re.Kode + ext

		fmt.Println("Start to scrap, store file to", fileName, re.Kode, i)

		if err := svc.ScrapAllWithStartingPoint(criterion, re, fileName); err != nil {
			fmt.Println("Error to scrap", i, re.Kode)
		}
	}

	return nil
}

func (svc *service) ScrapAllWithStartingPoint(criterion model.Criterion, startingPoint model.PPWTEntity, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if startingPoint.Tingkat != 0 && startingPoint.Parent == nil {
		startingPoint, err = svc.kpuRepo.GetPPWTParent(startingPoint)
		if err != nil {
			return err
		}
	}

	fmt.Println("Getting ppwt list for", startingPoint.GetCanonicalCode(), startingPoint.GetCanonicalName())

	ppwtList, err := svc.ScrapPPWTWithStartingPoint(startingPoint)
	if err != nil {
		return err
	}

	fmt.Println("Found ppwt list", len(ppwtList))

	hhcwListCh := make(chan model.HHCWEntity, len(ppwtList))

	go func() {
		wg := new(sync.WaitGroup)
		sm := semaphore.NewSemaphore(svc.maximumRunningThread)

		wg.Add(len(ppwtList))

		for _, entity := range ppwtList {

			sm.Acquire()
			go func(entity model.PPWTEntity) {
				defer sm.Release()
				defer wg.Done()

				res, err := svc.kpuRepo.GetHHCWInfo(entity)
				if err != nil {
					return
				}

				hhcwListCh <- res
			}(entity)
		}

		wg.Wait()
		close(hhcwListCh)
	}()

	w := csv.NewWriter(file)
	defer w.Flush()
	header := []string{
		"Kode",
		"Provinsi",
		"Kabupaten",
		"Kecamatan",
		"Kelurahan",
		"TPS",
		"Suara Paslon 01",
		"Suara Paslon 02",
		"Suara Paslon 03",
		"Jumlah",
		"Suara Sah",
		"Suara Tidak Sah",
		"Jumlah Suara",
		"DPT", "DPTb",
		"DPTk",
		"Jumlah Hak Pilih",
		"Selisih Suara Semua Paslon dengan Jumlah Suara Sah",
		"Selisih Suara Sah dan Tidak Sah dengan Total Suara",
		"Selisih Hak Pilih dengan Jumlah Suara",
		"Link Web KPU",
	}
	_ = w.Write(header)

	for {
		select {
		case hhwc, ok := <-hhcwListCh:
			if !ok {
				return nil
			}

			if criterion.IsMatchFor(hhwc) {
				ppwt := *hhwc.Parent
				link, err := svc.kpuRepo.GetPageLink(*hhwc.Parent)
				if err != nil {
					return err
				}
				fmt.Printf("%v %v\t%v\t%v\n", strings.Join(ppwt.GetCanonicalName(), "/"), ppwt.Kode, hhwc.Chart.String(), link)

				col := make([]string, 0, len(header))
				col = append(col, ppwt.Kode)
				col = append(col, ppwt.GetCanonicalName()...)
				col = append(col, strconv.Itoa(hhwc.Chart.Paslon01), strconv.Itoa(hhwc.Chart.Paslon02), strconv.Itoa(hhwc.Chart.Paslon03), strconv.Itoa(hhwc.Chart.Sum()))
				col = append(col, strconv.Itoa(hhwc.Administrasi.Suara.Sah), strconv.Itoa(hhwc.Administrasi.Suara.TidakSah), strconv.Itoa(hhwc.Administrasi.Suara.Total))
				col = append(col, strconv.Itoa(hhwc.Administrasi.PemilihDpt.Jumlah), strconv.Itoa(hhwc.Administrasi.PenggunaDptb.Jumlah), strconv.Itoa(hhwc.Administrasi.PenggunaNonDpt.Jumlah), strconv.Itoa(hhwc.Administrasi.PenggunaTotal.Jumlah))
				col = append(col, strconv.Itoa(int(math.Abs(float64(hhwc.Chart.Sum()-hhwc.Administrasi.Suara.Sah)))))
				col = append(col, strconv.Itoa(int(math.Abs(float64((hhwc.Administrasi.Suara.Sah+hhwc.Administrasi.Suara.TidakSah)-hhwc.Administrasi.Suara.Total)))))
				col = append(col, strconv.Itoa(int(math.Abs(float64(hhwc.Administrasi.PenggunaTotal.Jumlah-hhwc.Administrasi.Suara.Total)))))
				col = append(col, link)
				_ = w.Write(col)
			}
		}
	}
}
