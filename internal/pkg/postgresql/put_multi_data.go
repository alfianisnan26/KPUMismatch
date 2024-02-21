package postgresql

import (
	"fmt"
	"github.com/lib/pq"
	"kawalrealcount/internal/data/model"
	"strings"
)

const queryMultiData = `INSERT INTO %s (code, total_votes_01, total_votes_02, total_votes_03,
               total_sum_votes, total_valid_votes, total_invalid_votes, total_votes, dpt, dptb, dptk, jml_hak_pilih,
               selisih_suara_paslon_dan_jumlah_sah, selisih_suara_sah_tidak_sah_dan_total,
               link, pic_urls, updated_at, obtained_at, update_id, all_in)
VALUES %s
ON CONFLICT (code)
    DO UPDATE SET
                  total_votes_01 = EXCLUDED.total_votes_01,
                  total_votes_02 = EXCLUDED.total_votes_02,
                  total_votes_03 = EXCLUDED.total_votes_03,
                  total_sum_votes = EXCLUDED.total_sum_votes,
                  total_valid_votes = EXCLUDED.total_valid_votes,
                  total_invalid_votes = EXCLUDED.total_invalid_votes,
                  total_votes = EXCLUDED.total_votes,
                  dpt = EXCLUDED.dpt,
                  dptb = EXCLUDED.dptb,
                  dptk = EXCLUDED.dptk,
                  jml_hak_pilih = EXCLUDED.jml_hak_pilih,
                  selisih_suara_paslon_dan_jumlah_sah = EXCLUDED.selisih_suara_paslon_dan_jumlah_sah,
                  selisih_suara_sah_tidak_sah_dan_total = EXCLUDED.selisih_suara_sah_tidak_sah_dan_total,
                  link = EXCLUDED.link,
        		  pic_urls = EXCLUDED.pic_urls,
                  updated_at = EXCLUDED.updated_at,
                  obtained_at = EXCLUDED.obtained_at,
                  update_id = EXCLUDED.update_id,
                  all_in = EXCLUDED.all_in;`

const queryMultiDataPpwt = `INSERT INTO %s (code, provinsi, kabupaten, kecamatan,
               tps)
VALUES %s
ON CONFLICT (code)
    DO UPDATE SET
                	provinsi = EXCLUDED.provinsi,
    				kabupaten = EXCLUDED.kabupaten,
    				kecamatan = EXCLUDED.kecamatan,
    				tps = EXCLUDED.tps;`

func (r *repo) PutReplaceListData(entities []*model.HHCWEntity, updateId uint64) error {
	var args = make([]interface{}, 0, len(entities)*20)
	for _, entity := range entities {
		args = append(args, buildArgs(entity, updateId)...)
	}

	query := fmt.Sprintf(queryMultiData, r.tableRecord, buildPlaceholder(20, len(entities)))

	_, err := r.db.Exec(query,
		args...,
	)

	return err
}

func (r *repo) PutReplacePPWT(entities map[string]*model.PPWTEntity) error {
	var args = make([]interface{}, 0, len(entities)*5)
	for _, entity := range entities {
		switch {
		case entity.Parent == nil:
			fallthrough
		case entity.Parent.Parent == nil:
			fallthrough
		case entity.Parent.Parent.Parent == nil:
			continue
		}

		args = append(args, buildArgsPpwt(entity)...)
	}

	query := fmt.Sprintf(queryMultiDataPpwt, r.tableRecord, buildPlaceholder(5, len(entities)))

	_, err := r.db.Exec(query,
		args...,
	)

	return err
}

func buildPlaceholder(placeholder int, group int) string {
	var groupStr = make([]string, group)
	placeHolderStr := make([]string, placeholder)
	for i := 0; i < group; i++ {
		for j := 0; j < placeholder; j++ {
			placeHolderStr[j] = fmt.Sprintf("$%d", (j+1)+((placeholder)*i))
		}
		groupStr[i] = fmt.Sprintf("(%s)", strings.Join(placeHolderStr, ","))
	}

	return strings.Join(groupStr, ",")
}

func buildArgs(entity *model.HHCWEntity, updateId uint64) []interface{} {

	metric := entity.Evaluate()
	return []interface{}{
		entity.Code,                               // 1
		entity.Chart.Paslon01,                     // 7
		entity.Chart.Paslon02,                     // 8
		entity.Chart.Paslon03,                     // 9
		entity.Chart.Sum(),                        //10
		entity.Administrasi.Suara.Sah,             //11
		entity.Administrasi.Suara.TidakSah,        //12
		entity.Administrasi.Suara.Total,           //13
		entity.Administrasi.PenggunaDpt.Jumlah,    //14
		entity.Administrasi.PenggunaDptb.Jumlah,   //15
		entity.Administrasi.PenggunaNonDpt.Jumlah, //16
		entity.Administrasi.PenggunaTotal.Jumlah,  //17
		metric.DivChartSumSuaraSah,                //18
		metric.DivSahTidakSahTotal,                //19
		entity.Link,                               //20
		pq.Array(entity.Images),                   //21
		entity.UpdatedAt.UTC().UnixMilli(),        //22
		entity.ObtainedAt.UTC().UnixMilli(),       //23
		updateId,                                  //24
		entity.Chart.GetAllInPaslon(),             //25
	}
}

func buildArgsPpwt(entity *model.PPWTEntity) []interface{} {
	return []interface{}{
		entity.Kode,
		entity.Parent.Parent.Parent.Nama,
		entity.Parent.Parent.Nama,
		entity.Parent.Nama,
		entity.Nama,
	}
}
