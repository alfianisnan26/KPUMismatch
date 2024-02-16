package model

import (
	"sync"
	"time"
)

type Stats struct {
	// auto create
	ID        uint64
	CreatedAt time.Time

	// auto update
	Chart        ChartInfo
	Administrasi AdministrasiInfo

	MostUpdated  time.Time
	LeastUpdated time.Time

	CountMetric   Metric
	SumMetric     Metric
	HighestMetric Metric

	TopDivChartSumSuaraSah   string
	TopDivSahTidakSahTotal   string
	TopDivSuaraPenggunaTotal string

	TotalRecord int

	Progress           float32
	EstimateTime       time.Duration
	LastProgressUpdate time.Time

	// final value
	ProcessingTime time.Duration
	FinishedAt     time.Time

	mutex sync.Mutex
}

func (s *Stats) Update(progress float32, estTime time.Duration, startTime time.Time) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Progress = progress
	s.EstimateTime = estTime
	s.LastProgressUpdate = time.Now()
	s.ProcessingTime = s.LastProgressUpdate.Sub(startTime)
}

func (s *Stats) Finalize(startTime time.Time) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.FinishedAt = time.Now()
	s.Progress = 100
	s.EstimateTime = 0
	s.ProcessingTime = s.FinishedAt.Sub(startTime)
}

func (s *Stats) Evaluate(entity HHCWEntity) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.TotalRecord++

	s.Chart = s.Chart.Add(entity.Chart)
	s.Administrasi.Suara = s.Administrasi.Suara.Add(entity.Administrasi.Suara)
	s.Administrasi.PenggunaTotal = s.Administrasi.PenggunaTotal.Add(entity.Administrasi.PenggunaTotal)
	s.Administrasi.PemilihDpt = s.Administrasi.PemilihDpt.Add(entity.Administrasi.PemilihDpt)
	s.Administrasi.PenggunaDpt = s.Administrasi.PenggunaDpt.Add(entity.Administrasi.PenggunaDpt)
	s.Administrasi.PenggunaDptb = s.Administrasi.PenggunaDptb.Add(entity.Administrasi.PenggunaDptb)
	s.Administrasi.PenggunaNonDpt = s.Administrasi.PenggunaNonDpt.Add(entity.Administrasi.PenggunaNonDpt)
	if v := entity.UpdatedAt; v.After(s.MostUpdated) {
		s.MostUpdated = v
	}
	if v := entity.UpdatedAt; s.LeastUpdated.IsZero() || v.Before(s.LeastUpdated) {
		s.LeastUpdated = v
	}

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

	if s.HighestMetric.DivSuaraPenggunaTotal == m.DivSuaraPenggunaTotal {
		s.TopDivSuaraPenggunaTotal = entity.Parent.Kode
	}

}
