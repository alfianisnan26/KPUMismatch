package postgresql

import (
	"fmt"
	"kawalrealcount/internal/data/model"
)

func (r *repo) UpdateStats(stats *model.Stats) error {
	query := `UPDATE %s
				SET
					total_votes_01 = $1,
					total_votes_02 = $2,
					total_votes_03 = $3,
					total_sum_votes = $4,
					total_valid_votes = $5,
					total_invalid_votes = $6,
					total_votes = $7,
					dpt = $8,
					dptb = $9,
					dptk = $10,
					jml_hak_pilih = $11,
					count_div_chart_sum_suara_sah = $12,
					count_div_sah_tidak_sah_total = $13,
					sum_div_chart_sum_suara_sah = $14,
					sum_div_sah_tidak_sah_total = $15,
					highest_div_chart_sum_suara_sah = $16,
					highest_div_sah_tidak_sah_total = $17,
					top_div_chart_sum_suara_sah = $18,
					top_div_sah_tidak_sah_total = $19,
					total_record = $20,
					progress = $21,
					estimate_time = $22,
					last_progress_update = $23,
					processing_time = $24,
					finished_at = $25,
					total_all_in_01 = $26,
					total_all_in_02 = $27,
					total_all_in_03 = $28,
				    total_clear_votes_01 = $29,
				    total_clear_votes_02 = $30,
				    total_clear_votes_03 = $31,
				    total_non_null_record = $32,
				    total_valid_non_null_record = $33
				WHERE id = $34;`

	var finishedAt int64
	if !stats.FinishedAt.IsZero() {
		finishedAt = stats.FinishedAt.UTC().UnixMilli()
	}

	_, err := r.db.Exec(fmt.Sprintf(query, r.tableStat),
		stats.Chart.Paslon01,                       //1
		stats.Chart.Paslon02,                       //2
		stats.Chart.Paslon03,                       //3
		stats.Chart.Sum(),                          //4
		stats.Administrasi.Suara.Sah,               //5
		stats.Administrasi.Suara.TidakSah,          //6
		stats.Administrasi.Suara.Total,             //7
		stats.Administrasi.PemilihDpt.Jumlah,       //8
		stats.Administrasi.PenggunaDptb.Jumlah,     //9
		stats.Administrasi.PenggunaNonDpt.Jumlah,   //10
		stats.Administrasi.PenggunaTotal.Jumlah,    //11
		stats.CountMetric.DivChartSumSuaraSah,      //12
		stats.CountMetric.DivSahTidakSahTotal,      //13
		stats.SumMetric.DivChartSumSuaraSah,        //14
		stats.SumMetric.DivSahTidakSahTotal,        //15
		stats.HighestMetric.DivChartSumSuaraSah,    //16
		stats.HighestMetric.DivSahTidakSahTotal,    //17
		stats.TopDivChartSumSuaraSah,               //18
		stats.TopDivSahTidakSahTotal,               //19
		stats.TotalRecord,                          //20
		stats.Progress,                             //21
		stats.EstimateTime.Milliseconds(),          //22
		stats.LastProgressUpdate.UTC().UnixMilli(), //23
		stats.ProcessingTime.Milliseconds(),        //24
		finishedAt,                                 //25
		stats.AllInChart.Paslon01,                  //26
		stats.AllInChart.Paslon02,                  //27
		stats.AllInChart.Paslon03,                  //28
		stats.ClearChart.Paslon01,                  //29
		stats.ClearChart.Paslon02,                  //30
		stats.ClearChart.Paslon03,                  //31
		stats.TotalNonNullRecord,                   //32
		stats.TotalValidNonNullRecord,              //33
		stats.ID,                                   //34
	)

	return err
}
