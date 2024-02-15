package scrapper

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"kawalrealcount/internal/data/model"
	"math"
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
	SheetSelisihJumlahDataSuaraSah        = "Selisih Jumlah Data Suara Sah"
	SheetSelisihJumlahSuaraDenganTotal    = "Selisih Jumlah Suara dengan Total"
	SheetSelisihJumlahHakPilihDenganSuara = "Selisih Jumlah Hak Pilih dengan Suara"
	SheetTPSAllIn                         = "TPS All In"
)

var sheet = []string{
	SheetSelisihJumlahDataSuaraSah,
	SheetSelisihJumlahSuaraDenganTotal,
	SheetSelisihJumlahHakPilihDenganSuara,
	SheetTPSAllIn,
}

func (svc *service) ScrapAllCompiled(filePath string) error {
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
		fmt.Println("Start to scrap, store file to", filePath, re.Nama, sheetMap)

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

	LoopFor:
		for {
			select {
			case hhwc, ok := <-hhcwListCh:
				if !ok {
					break LoopFor
				}

				// filter selisih data paslon
				if sum, sah := hhwc.Chart.Sum(), hhwc.Administrasi.Suara.Sah; sum != 0 && sah != 0 && sum != sah {
					if err := svc.writeCell(f, &sheetMap, SheetSelisihJumlahDataSuaraSah, hhwc); err != nil {
						return err
					}
				}

				// filter selisih total
				if s, ts, tot := hhwc.Administrasi.Suara.Sah, hhwc.Administrasi.Suara.TidakSah, hhwc.Administrasi.Suara.Total; s+ts != 0 && tot != 0 && s+ts != tot {
					if err := svc.writeCell(f, &sheetMap, SheetSelisihJumlahSuaraDenganTotal, hhwc); err != nil {
						return err
					}
				}
				// filter selisih hak pilih
				if st, pt := hhwc.Administrasi.Suara.Total, hhwc.Administrasi.PenggunaTotal.Jumlah; st != pt {
					if err := svc.writeCell(f, &sheetMap, SheetSelisihJumlahHakPilihDenganSuara, hhwc); err != nil {
						return err
					}
				}

				// filter all in
				if hhwc.Chart.IsAllIn() {
					if err := svc.writeCell(f, &sheetMap, SheetTPSAllIn, hhwc); err != nil {
						return err
					}
				}
			default:
			}

		}
	}

	for _, v := range f.GetSheetMap() {
		if _, found := sheetMap[v]; !found {
			f.DeleteSheet(v)
		}
	}
	return f.SaveAs(filePath)
}

func (svc *service) writeCell(f *excelize.File, sheetmap *map[string]int, sheet string, data model.HHCWEntity) error {

	link, err := svc.kpuRepo.GetPageLink(*data.Parent)
	if err != nil {
		return err
	}

	var row = make([]interface{}, 0, len(header))

	row = append(row, data.UpdatedAt.Format(time.DateTime), data.Parent.Kode)
	row = append(row, data.Parent.GetCanonicalName()...)
	row = append(row, data.Chart.Paslon01, data.Chart.Paslon02, data.Chart.Paslon03, data.Chart.Sum())
	row = append(row, data.Administrasi.Suara.Sah, data.Administrasi.Suara.TidakSah, data.Administrasi.Suara.Total)
	row = append(row, data.Administrasi.PemilihDpt.Jumlah, data.Administrasi.PenggunaDptb.Jumlah, data.Administrasi.PenggunaNonDpt.Jumlah, data.Administrasi.PenggunaTotal.Jumlah)
	row = append(row, int(math.Abs(float64(data.Chart.Sum()-data.Administrasi.Suara.Sah))))
	row = append(row, int(math.Abs(float64((data.Administrasi.Suara.Sah+data.Administrasi.Suara.TidakSah)-data.Administrasi.Suara.Total))))
	row = append(row, int(math.Abs(float64(data.Administrasi.Suara.Total-data.Administrasi.PenggunaTotal.Jumlah))))
	row = append(row, link)

	if (*sheetmap)[sheet]%100 == 0 {
		fmt.Printf("Sample %% 100: %s\t| %s\n", data.String(), link)
	}

	for col, value := range row {
		cell := excelize.ToAlphaString(col) + strconv.Itoa((*sheetmap)[sheet]) // Write in first row
		f.SetCellValue(sheet, cell, value)
	}

	(*sheetmap)[sheet]++

	return nil
}
