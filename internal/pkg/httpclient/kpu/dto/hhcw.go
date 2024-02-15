package dto

import (
	"kawalrealcount/internal/data/model"
	"time"
)

type HHCWEntity struct {
	Chart struct {
		Paslon01 int         `json:"100025"`
		Paslon02 int         `json:"100026"`
		Paslon03 int         `json:"100027"`
		Null     interface{} `json:"null"`
	} `json:"chart"`
	Images       []string `json:"images"`
	Administrasi struct {
		SuaraSah        int `json:"suara_sah"`
		SuaraTotal      int `json:"suara_total"`
		PemilihDptJ     int `json:"pemilih_dpt_j"`
		PemilihDptL     int `json:"pemilih_dpt_l"`
		PemilihDptP     int `json:"pemilih_dpt_p"`
		PenggunaDptJ    int `json:"pengguna_dpt_j"`
		PenggunaDptL    int `json:"pengguna_dpt_l"`
		PenggunaDptP    int `json:"pengguna_dpt_p"`
		PenggunaDptbJ   int `json:"pengguna_dptb_j"`
		PenggunaDptbL   int `json:"pengguna_dptb_l"`
		PenggunaDptbP   int `json:"pengguna_dptb_p"`
		SuaraTidakSah   int `json:"suara_tidak_sah"`
		PenggunaTotalJ  int `json:"pengguna_total_j"`
		PenggunaTotalL  int `json:"pengguna_total_l"`
		PenggunaTotalP  int `json:"pengguna_total_p"`
		PenggunaNonDptJ int `json:"pengguna_non_dpt_j"`
		PenggunaNonDptL int `json:"pengguna_non_dpt_l"`
		PenggunaNonDptP int `json:"pengguna_non_dpt_p"`
	} `json:"administrasi"`
	Psu         interface{} `json:"psu"`
	Ts          string      `json:"ts"`
	StatusSuara bool        `json:"status_suara"`
	StatusAdm   bool        `json:"status_adm"`
}

func (ti HHCWEntity) ToModel(ppwt model.PPWTEntity) (model.HHCWEntity, error) {
	updatedAt, err := time.Parse(time.DateTime, ti.Ts)
	if err != nil {
		return model.HHCWEntity{}, err
	}

	return model.HHCWEntity{
		Parent: &ppwt,
		Chart: model.ChartInfo{
			Paslon01: ti.Chart.Paslon01,
			Paslon02: ti.Chart.Paslon02,
			Paslon03: ti.Chart.Paslon03,
		},
		Images: ti.Images,
		Administrasi: model.AdministrasiInfo{
			Suara: model.SuaraData{
				Sah:      ti.Administrasi.SuaraSah,
				TidakSah: ti.Administrasi.SuaraTidakSah,
				Total:    ti.Administrasi.SuaraTotal,
			},
			PemilihDpt: model.JLPData{
				Jumlah:    ti.Administrasi.PemilihDptJ,
				LakiLaki:  ti.Administrasi.PemilihDptL,
				Perempuan: ti.Administrasi.PemilihDptP,
			},
			PenggunaDpt: model.JLPData{
				Jumlah:    ti.Administrasi.PenggunaDptJ,
				LakiLaki:  ti.Administrasi.PenggunaDptL,
				Perempuan: ti.Administrasi.PenggunaDptP,
			},
			PenggunaDptb: model.JLPData{
				Jumlah:    ti.Administrasi.PenggunaDptbJ,
				LakiLaki:  ti.Administrasi.PenggunaDptbL,
				Perempuan: ti.Administrasi.PenggunaDptbP,
			},
			PenggunaTotal: model.JLPData{
				Jumlah:    ti.Administrasi.PenggunaTotalJ,
				LakiLaki:  ti.Administrasi.PenggunaTotalL,
				Perempuan: ti.Administrasi.PenggunaTotalP,
			},
			PenggunaNonDpt: model.JLPData{
				Jumlah:    ti.Administrasi.PenggunaNonDptJ,
				LakiLaki:  ti.Administrasi.PenggunaNonDptL,
				Perempuan: ti.Administrasi.PenggunaNonDptP,
			},
		},
		UpdatedAt:          updatedAt,
		ObtainedAt:         time.Now(),
		StatusSuara:        ti.StatusSuara,
		StatusAdministrasi: ti.StatusAdm,
	}, nil
}
