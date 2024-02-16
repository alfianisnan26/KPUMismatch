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
					count_div_suara_pengguna_total = $14,
					sum_div_chart_sum_suara_sah = $15,
					sum_div_sah_tidak_sah_total = $16,
					sum_div_suara_pengguna_total = $17,
					highest_div_chart_sum_suara_sah = $18,
					highest_div_sah_tidak_sah_total = $19,
					highest_div_suara_pengguna_total = $20,
					top_div_chart_sum_suara_sah = $21,
					top_div_sah_tidak_sah_total = $22,
					top_div_suara_pengguna_total = $23,
					total_record = $24,
					progress = $25,
					estimate_time = $26,
					last_progress_update = $27,
					processing_time = $28,
					finished_at = $29
				WHERE id = $30;`

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
		stats.CountMetric.DivSuaraPenggunaTotal,    //14
		stats.SumMetric.DivChartSumSuaraSah,        //15
		stats.SumMetric.DivSahTidakSahTotal,        //16
		stats.SumMetric.DivSuaraPenggunaTotal,      //17
		stats.HighestMetric.DivChartSumSuaraSah,    //18
		stats.HighestMetric.DivSahTidakSahTotal,    //19
		stats.HighestMetric.DivSuaraPenggunaTotal,  //20
		stats.TopDivChartSumSuaraSah,               //21
		stats.TopDivSahTidakSahTotal,               //22
		stats.TopDivSuaraPenggunaTotal,             //23
		stats.TotalRecord,                          //24
		stats.Progress,                             //25
		stats.EstimateTime.Milliseconds(),          //26
		stats.LastProgressUpdate.UTC().UnixMilli(), //27
		stats.ProcessingTime.Milliseconds(),        //28
		finishedAt,                                 //29
		stats.ID,                                   //30
	)

	return err
}
