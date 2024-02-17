package model

type Metric struct {
	DivChartSumSuaraSah int `json:"div_chart_sum_suara_sah,omitempty"`
	DivSahTidakSahTotal int `json:"div_sah_tidak_sah_total,omitempty"`
}

func (m *Metric) Count() Metric {
	var o Metric

	if m.DivChartSumSuaraSah > 0 {
		o.DivChartSumSuaraSah = 1
	}

	if m.DivSahTidakSahTotal > 0 {
		o.DivSahTidakSahTotal = 1
	}

	return o
}

func (m *Metric) Add(o Metric) Metric {
	return Metric{
		DivChartSumSuaraSah: m.DivChartSumSuaraSah + o.DivChartSumSuaraSah,
		DivSahTidakSahTotal: m.DivSahTidakSahTotal + o.DivSahTidakSahTotal,
	}
}

func (m *Metric) Max(o Metric) Metric {
	return Metric{
		DivChartSumSuaraSah: max(m.DivChartSumSuaraSah, o.DivChartSumSuaraSah),
		DivSahTidakSahTotal: max(m.DivSahTidakSahTotal, o.DivSahTidakSahTotal),
	}
}
