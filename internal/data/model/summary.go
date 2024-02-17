package model

import (
	"encoding/json"
	"fmt"
)

type SummaryModule struct {
	Chart     ChartInfo `json:"chart,omitempty"`
	Suara     SuaraData `json:"suara,omitempty"`
	HakPilih  int       `json:"hak_pilih,omitempty"`
	SumMetric Metric    `json:"sum_metric,omitempty"`
	TotalData int       `json:"total_data,omitempty"`
}

func (m SummaryModule) RawJson() []byte {
	rawJson, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	return rawJson
}

type Summary struct {
	RawData     SummaryModule
	NotNullData SummaryModule
	ClearData   SummaryModule
	AllInData   SummaryModule
}
