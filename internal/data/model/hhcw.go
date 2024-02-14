package model

import (
	"fmt"
	"sort"
	"time"
)

type HHCWEntity struct {
	Chart              ChartInfo
	Images             []string
	Administrasi       AdministrasiInfo
	UpdatedAt          time.Time
	StatusSuara        bool
	StatusAdministrasi bool

	Parent *PPWTEntity
}

func (hhcw HHCWEntity) String() string {
	parent := ""

	if hhcw.Parent != nil {
		parent = (*hhcw.Parent).Kode
	}

	return fmt.Sprintf("(%v) %v > %v \t | Suara: %v \t | Pengguna: %v", hhcw.UpdatedAt, parent, hhcw.Chart.String(), hhcw.Administrasi.Suara.String(), hhcw.Administrasi.PenggunaTotal.String())
}

type ChartInfo struct {
	Paslon01,
	Paslon02,
	Paslon03 int

	highestPercentage *float32
}

func (ci *ChartInfo) String() string {
	return fmt.Sprintf("01:%d | 02:%d | 03:%d | Sum:%d | Dist:%.2f%%", ci.Paslon01, ci.Paslon02, ci.Paslon03, ci.Sum(), ci.GetHighestDeltaPercentage())
}

func (ci *ChartInfo) Sum() int {
	return ci.Paslon01 + ci.Paslon02 + ci.Paslon03
}

func (ci *ChartInfo) GetHighestDeltaPercentage() float32 {
	if ci.highestPercentage != nil {
		return *ci.highestPercentage
	}

	num := []int{ci.Paslon01, ci.Paslon02, ci.Paslon03, 0}
	sort.Ints(num)
	d2 := num[3] - num[2]
	d1 := num[2] - num[1]
	d0 := num[1] - num[0]

	total := d0 + d1 + d2
	highest := max(d0, d1, d2)
	percentage := float32(highest) / float32(total) * 100
	ci.highestPercentage = &percentage
	return *ci.highestPercentage
}

type AdministrasiInfo struct {
	Suara SuaraData
	PemilihDpt,
	PenggunaDpt,
	PenggunaDptb,
	PenggunaTotal,
	PenggunaNonDpt JLPData
}

type JLPData struct {
	Jumlah,
	LakiLaki,
	Perempuan int
}

func (jd JLPData) String() string {
	return fmt.Sprintf("L:%d | P:%d | Sum:%d | Valid:%v", jd.LakiLaki, jd.Perempuan, jd.Jumlah, jd.IsValid())
}

func (jd JLPData) IsValid() bool {
	return jd.LakiLaki+jd.Perempuan == jd.Jumlah
}

type SuaraData struct {
	Sah,
	TidakSah,
	Total int
}

func (sd SuaraData) String() string {
	return fmt.Sprintf("S:%d | TS:%d | Sum:%d | Valid:%v", sd.Sah, sd.TidakSah, sd.Total, sd.IsValid())
}

func (sd SuaraData) IsValid() bool {
	return sd.Sah+sd.TidakSah == sd.Total
}
