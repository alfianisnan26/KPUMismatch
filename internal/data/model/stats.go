package model

import (
	"time"
)

type Stats struct {
	// auto create
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

	TotalNonNullRecord      int
	TotalValidNonNullRecord int

	WebStast *WebStats

	// final value
	FinishedAt time.Time
}

func (s *Stats) Finalize() {
	s.FinishedAt = s.WebStast.Timestamp
	s.WebStast.Percentage = 100
	s.WebStast.Estimation = 0
}

func (s *Stats) Evaluate() {
	s.WebStast.DataCount++
}
