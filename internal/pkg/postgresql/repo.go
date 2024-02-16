package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"kawalrealcount/internal/data/dao"
	"kawalrealcount/internal/data/model"
)

type repo struct {
	db          *sql.DB
	tableRecord string
	tableStat   string
}

func (r *repo) PutReplaceData(entity model.HHCWEntity, updateId uint64) error {
	canonical := entity.Parent.GetCanonicalName()
	query := `INSERT INTO %s (code, provinsi, kabupaten, kecamatan, kelurahan, tps, total_votes_01, total_votes_02, total_votes_03,
               total_sum_votes, total_valid_votes, total_invalid_votes, total_votes, dpt, dptb, dptk, jml_hak_pilih,
               selisih_suara_paslon_dan_jumlah_sah, selisih_suara_sah_tidak_sah_dan_total,
               link, pic_urls, updated_at, obtained_at, update_id, all_in)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
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

	metric := entity.Evaluate()
	_, err := r.db.Exec(fmt.Sprintf(query, r.tableRecord),
		entity.Parent.Kode,                        // 1
		canonical[0],                              // 2
		canonical[1],                              // 3
		canonical[2],                              // 4
		canonical[3],                              // 5
		canonical[4],                              // 6
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
	)

	return err
}

type Param struct {
	ConnectionURL string
	TableRecord   string
	TableStats    string
}

func New(param Param) (dao.Database, error) {
	// Initialize PostgreSQL connection
	db, err := sql.Open("postgres", param.ConnectionURL)
	if err != nil {
		return nil, err
	}

	// Check if the connection is successful
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Return the PostgreSQL repository
	return &repo{
		db:          db,
		tableRecord: param.TableRecord,
		tableStat:   param.TableStats,
	}, nil
}
