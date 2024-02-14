package model

type Criterion struct {
	DeltaErrThreshold float32

	WithDeltaErrThreshold               bool
	WithInvalidSumOfSuara               bool
	WithInvalidSumOfPengguna            bool
	WithNonZeroMismatchSuaraAndPengguna bool
	WithNonZeroMismatchSuara            bool

	IgnoreAll bool
}

func (c Criterion) IsMatchFor(hhcw HHCWEntity) bool {
	sum := hhcw.Chart.Sum()
	switch {
	case c.IgnoreAll:
		fallthrough
	case c.WithDeltaErrThreshold && hhcw.Chart.GetHighestDeltaPercentage() > c.DeltaErrThreshold:
		fallthrough
	case c.WithInvalidSumOfSuara && !hhcw.Administrasi.Suara.IsValid():
		fallthrough
	case c.WithInvalidSumOfPengguna && !hhcw.Administrasi.PenggunaTotal.IsValid():
		fallthrough
	case c.WithNonZeroMismatchSuaraAndPengguna && hhcw.Administrasi.Suara.Total != 0 && hhcw.Administrasi.PenggunaTotal.Jumlah != 0 && hhcw.Administrasi.Suara.Total != hhcw.Administrasi.PenggunaTotal.Jumlah:
		fallthrough
	case c.WithNonZeroMismatchSuara && sum != 0 && hhcw.Administrasi.Suara.Total != 0 && hhcw.Administrasi.Suara.Sah != sum:
		return true
	default:
		return false
	}
}
