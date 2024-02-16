package model

import (
	"sync"
	"time"
)

type Stats struct {
	// auto create
	ID          uint64
	CreatedAt   time.Time
	Contributor string

	// auto update
	Chart      ChartInfo
	ClearChart ChartInfo
	AllInChart ChartInfo

	Administrasi AdministrasiInfo

	CountMetric   Metric
	SumMetric     Metric
	HighestMetric Metric

	TopDivChartSumSuaraSah string
	TopDivSahTidakSahTotal string

	TotalRecord             int
	TotalNonNullRecord      int
	TotalValidNonNullRecord int

	Progress           float32
	EstimateTime       time.Duration
	LastProgressUpdate time.Time

	// final value
	ProcessingTime time.Duration
	FinishedAt     time.Time

	mutex sync.Mutex
}

func (s *Stats) Update(progress float32, estTime time.Duration, startTime time.Time, count int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Progress = progress
	s.EstimateTime = estTime
	s.LastProgressUpdate = time.Now()
	s.ProcessingTime = s.LastProgressUpdate.Sub(startTime)
	s.TotalRecord = count
}

func (s *Stats) Finalize(startTime time.Time, count int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.FinishedAt = time.Now()
	s.Progress = 100
	s.EstimateTime = 0
	s.ProcessingTime = s.FinishedAt.Sub(startTime)
	s.TotalRecord = count
}

func (s *Stats) Evaluate(entity HHCWEntity) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Chart = s.Chart.Add(entity.Chart)
	var allInChart ChartInfo

	allIn := entity.Chart.GetAllInPaslon()
	switch allIn {
	case 1:
		allInChart.Paslon01 = 1
	case 2:
		allInChart.Paslon02 = 1
	case 3:
		allInChart.Paslon03 = 1
	default:
	}

	s.AllInChart = s.AllInChart.Add(allInChart)
	s.Administrasi.Suara = s.Administrasi.Suara.Add(entity.Administrasi.Suara)
	s.Administrasi.PenggunaTotal = s.Administrasi.PenggunaTotal.Add(entity.Administrasi.PenggunaTotal)
	s.Administrasi.PemilihDpt = s.Administrasi.PemilihDpt.Add(entity.Administrasi.PemilihDpt)
	s.Administrasi.PenggunaDpt = s.Administrasi.PenggunaDpt.Add(entity.Administrasi.PenggunaDpt)
	s.Administrasi.PenggunaDptb = s.Administrasi.PenggunaDptb.Add(entity.Administrasi.PenggunaDptb)
	s.Administrasi.PenggunaNonDpt = s.Administrasi.PenggunaNonDpt.Add(entity.Administrasi.PenggunaNonDpt)

	m := entity.Evaluate()

	s.CountMetric = s.CountMetric.Add(m.Count())
	s.SumMetric = s.SumMetric.Add(m)
	s.HighestMetric = s.HighestMetric.Max(m)

	if s.HighestMetric.DivChartSumSuaraSah == m.DivChartSumSuaraSah {
		s.TopDivChartSumSuaraSah = entity.Parent.Kode
	}

	if s.HighestMetric.DivSahTidakSahTotal == m.DivSahTidakSahTotal {
		s.TopDivSahTidakSahTotal = entity.Parent.Kode
	}

	if entity.IsNonNullVote() {
		s.TotalNonNullRecord++

		if entity.IsValidVote() {
			s.TotalValidNonNullRecord++
			s.ClearChart = s.ClearChart.Add(entity.Chart)
		}
	}

}
