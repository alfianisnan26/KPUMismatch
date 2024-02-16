package model

type Metric struct {
	DivChartSumSuaraSah   int
	DivSahTidakSahTotal   int
	DivSuaraPenggunaTotal int
}

func (m *Metric) Count() Metric {
	var o Metric

	if m.DivSuaraPenggunaTotal > 0 {
		o.DivSuaraPenggunaTotal = 1
	}

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
		DivChartSumSuaraSah:   m.DivChartSumSuaraSah + o.DivChartSumSuaraSah,
		DivSahTidakSahTotal:   m.DivSahTidakSahTotal + o.DivSahTidakSahTotal,
		DivSuaraPenggunaTotal: m.DivSuaraPenggunaTotal + o.DivSuaraPenggunaTotal,
	}
}

func (m *Metric) Max(o Metric) Metric {
	return Metric{
		DivChartSumSuaraSah:   max(m.DivChartSumSuaraSah, o.DivChartSumSuaraSah),
		DivSahTidakSahTotal:   max(m.DivSahTidakSahTotal, o.DivSahTidakSahTotal),
		DivSuaraPenggunaTotal: max(m.DivSuaraPenggunaTotal, o.DivSuaraPenggunaTotal),
	}
}
