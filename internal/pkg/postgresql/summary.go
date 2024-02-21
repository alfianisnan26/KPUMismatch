package postgresql

import (
	"fmt"
	"kawalrealcount/internal/data/model"
	"time"
)

const (
	whereNonNull           = `WHERE total_sum_votes != 0 and total_votes != 0`
	whereClearValid        = `WHERE total_sum_votes != 0 and total_votes != 0 and selisih_suara_paslon_dan_jumlah_sah = 0 and selisih_suara_sah_tidak_sah_dan_total = 0`
	whereAllInValid        = `WHERE all_in > 0 and total_sum_votes != 0 and total_votes != 0 and selisih_suara_paslon_dan_jumlah_sah = 0 and selisih_suara_sah_tidak_sah_dan_total = 0`
	whereErrorBeyondMaxVal = `WHERE total_sum_votes > 300 OR total_votes > 300 OR jml_hak_pilih > 300`
)

const querySelectSum = `
SELECT
	sum(total_votes_01) as total_votes_01,
	sum(total_votes_02) as total_votes_02,
	sum(total_votes_03) as total_votes_03,
	sum(total_valid_votes) as total_valid_votes,
	sum(total_invalid_votes) as total_invalid_votes,
	sum(total_votes) as total_votes,
	sum(jml_hak_pilih) as jml_hak_pilih,
	sum(selisih_suara_paslon_dan_jumlah_sah) as selisih_suara_paslon_dan_jumlah_sah,
	sum(selisih_suara_sah_tidak_sah_dan_total) as selisih_suara_sah_tidak_sah_dan_total,
	count(*) as total`

func (r *repo) GetSummary() (model.Summary, error) {
	whereSeq := []string{"", whereNonNull, whereClearValid, whereAllInValid, whereErrorBeyondMaxVal}
	var summary model.Summary
	for i, where := range whereSeq {
		var v model.SummaryModule
		q := fmt.Sprintf(
			"%s FROM %s %s",
			querySelectSum,
			r.tableRecord,
			where,
		)
		if err := r.db.QueryRow(
			q,
		).Scan(
			&v.Chart.Paslon01,
			&v.Chart.Paslon02,
			&v.Chart.Paslon03,
			&v.Suara.Sah,
			&v.Suara.TidakSah,
			&v.Suara.Total,
			&v.HakPilih,
			&v.SumMetric.DivChartSumSuaraSah,
			&v.SumMetric.DivSahTidakSahTotal,
			&v.TotalData,
		); err != nil {
			return model.Summary{}, err
		}

		switch i {
		case 0:
			summary.RawData = v
		case 1:
			summary.NotNullData = v
		case 2:
			summary.ClearData = v
		case 3:
			summary.AllInData = v
		case 4:
			summary.ErrorBeyondMaxValData = v
		}
	}

	return summary, nil
}

func (r *repo) InsertSummary(summary model.Summary) error {

	ts := time.Now()

	query := `INSERT INTO %s(raw_data, not_null_data, clear_data, all_in_data, error_max_val_data, ts, ts_unix) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(fmt.Sprintf(query, r.tableHistogram),
		summary.RawData.RawJson(),
		summary.NotNullData.RawJson(),
		summary.ClearData.RawJson(),
		summary.AllInData.RawJson(),
		summary.ErrorBeyondMaxValData.RawJson(),
		ts,
		ts.UnixMilli(),
	)

	return err
}
