package dto

import "kawalrealcount/internal/data/model"

type PPWTEntity struct {
	ID      int64  `json:"id,omitempty"`
	Kode    string `json:"kode,omitempty"`
	Nama    string `json:"nama,omitempty"`
	Tingkat int    `json:"tingkat,omitempty"`
}

func (e PPWTEntity) ToModel(parent model.PPWTEntity) model.PPWTEntity {

	res := model.PPWTEntity{
		ID:      e.ID,
		Kode:    e.Kode,
		Nama:    e.Nama,
		Tingkat: e.Tingkat,
	}

	if parent.Kode != "0" {
		res.Parent = &parent
	}

	return res
}
