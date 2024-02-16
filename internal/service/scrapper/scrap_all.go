package scrapper

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"kawalrealcount/internal/data/model"
	"strconv"
	"time"
)

var header = []string{
	"Terakhir Diperbarui",
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

const (
	SheetSelisihData = "Selisih Data"
	SheetTPSAllIn    = "TPS All In"
)

var sheet = []string{
	SheetSelisihData,
	SheetTPSAllIn,
}

func (svc *service) ScrapAllCompiled(filePath string, isScrapAll bool) error {
	startingPoint := model.NewPPWT("0")

	res, err := svc.kpuRepo.GetPPWTList(startingPoint)
	if err != nil {
		return err
	}

	f := excelize.NewFile()

	sheetMap := make(map[string]int, len(sheet))

	for _, s := range sheet {
		_ = f.NewSheet(s)
		for col, value := range header {
			cell := excelize.ToAlphaString(col) + "1" // Write in first row
			f.SetCellValue(s, cell, value)
		}

		sheetMap[s] = 2
	}

	for i, re := range res {
		var percentage float32
		if i != 0 {
			percentage = float32(i) / float32(len(res)) * 100
		}

		fmt.Printf("[%03.2f%%] Start to scrap, store file to: %s | %s | %v\n", percentage, filePath, re.Nama, sheetMap)

		ppwtList, err := svc.ScrapPPWTWithStartingPoint(startingPoint)
		if err != nil {
			return err
		}

		length := len(ppwtList)

		fmt.Println("Found ppwt list", length)

		hhcwListCh := make(chan model.HHCWEntity, length)

		if err := svc.ScrapAllWithStartingPoint(hhcwListCh, re, ppwtList); err != nil {
			fmt.Println("Error to scrap", i, re.Kode)
		}

		var count int

	LoopFor:
		for {
			select {
			case hhwc, ok := <-hhcwListCh:
				if !ok {
					break LoopFor
				}
				count++
				subPercentageRaw := float32(count) / float32(length)
				subPercentage := subPercentageRaw * 100
				deltaPercentage := percentage + (((float32(i+1) / float32(len(res)) * 100) - percentage) * subPercentageRaw)

				link, err := svc.kpuRepo.GetPageLink(*hhwc.Parent)
				if err != nil {
					return err
				}
				hhwc.Link = link

				var store = isScrapAll
				// filter selisih data paslon
				s, ts, tot := hhwc.Administrasi.Suara.Sah, hhwc.Administrasi.Suara.TidakSah, hhwc.Administrasi.Suara.Total
				sum, sah := hhwc.Chart.Sum(), hhwc.Administrasi.Suara.Sah
				if (sum != 0 && sah != 0 && sum != sah) || (s+ts != 0 && tot != 0 && s+ts != tot) {
					store = true
					if filePath != "" {
						if err := svc.writeCell(f, &sheetMap, SheetSelisihData, hhwc); err != nil {
							return err
						}
					}
				}

				// filter all in
				if hhwc.Chart.IsAllIn() {
					store = true
					if filePath != "" {
						if err := svc.writeCell(f, &sheetMap, SheetTPSAllIn, hhwc); err != nil {
							return err
						}

					}
				}

				if store && svc.databaseRepo != nil {
					if count%100 == 0 {
						fmt.Printf("[%03.2f%%][%03.2f%%]\t%d | %s\t| %s\n", deltaPercentage, subPercentage, count, hhwc.String(), hhwc.Link)
					}
					if err := svc.databaseRepo.PutReplaceData(hhwc); err != nil {
						// ignore error
						fmt.Println(err.Error())
					}
				}

			default:
			}

		}
	}

	if filePath != "" {
		return f.SaveAs(filePath)
	}

	return nil
}

func (svc *service) writeCell(f *excelize.File, sheetmap *map[string]int, sheet string, data model.HHCWEntity) error {

	var row = make([]interface{}, 0, len(header))

	metric := data.Evaluate()

	row = append(row, data.UpdatedAt.Format(time.DateTime), data.Parent.Kode)
	row = append(row, data.Parent.GetCanonicalName()...)
	row = append(row, data.Chart.Paslon01, data.Chart.Paslon02, data.Chart.Paslon03, data.Chart.Sum())
	row = append(row, data.Administrasi.Suara.Sah, data.Administrasi.Suara.TidakSah, data.Administrasi.Suara.Total)
	row = append(row, data.Administrasi.PemilihDpt.Jumlah, data.Administrasi.PenggunaDptb.Jumlah, data.Administrasi.PenggunaNonDpt.Jumlah, data.Administrasi.PenggunaTotal.Jumlah)
	row = append(row, metric.DivChartSumSuaraSah)
	row = append(row, metric.DivSahTidakSahTotal)
	row = append(row, metric.DivSuaraPenggunaTotal)
	row = append(row, data.Link)

	for col, value := range row {
		cell := excelize.ToAlphaString(col) + strconv.Itoa((*sheetmap)[sheet]) // Write in first row
		f.SetCellValue(sheet, cell, value)
	}

	(*sheetmap)[sheet]++

	return nil
}
