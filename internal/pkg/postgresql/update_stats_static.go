package postgresql

import (
	"fmt"
	"kawalrealcount/internal/data/dao"
	"kawalrealcount/internal/data/model"
	"time"
)

var _ dao.UpdaterDatabase = &repo{}

const queryMapDist = `SELECT
CASE
WHEN "01" >= "02" AND "02" >= "03" THEN 'Paslon 01'
WHEN "02" >= "01" AND "02" >= "03" THEN 'Paslon 02'
ELSE 'Paslon 03'
END AS "Paslon Menang",
GREATEST("01", "02", "03") as votes_value,
"01" as "Paslon 01",
"02" as "Paslon 02",
"03" as "Paslon 03",
latitude,
longitude,
kabupaten as "Kabupaten/Kota"
FROM (
SELECT
sum(total_votes_01) as "01",
sum(total_votes_02) as "02",
sum(total_votes_03) as "03",
SUBSTRING(code, 0, 5) as city_id,
kabupaten
FROM
%s
GROUP BY
city_id, kabupaten
) as ci INNER JOIN map_loc as ml ON ml.serial = ci.city_id::numeric WHERE city_id ~ E'^\\d+$'`

const queryPotentialTableSum = `SELECT 
    code AS "Kode",
    provinsi AS "Provinsi",
    kabupaten AS "Kabupaten",
    kecamatan AS "Kecamatan",
    kelurahan AS "Kelurahan",
    tps AS "TPS",
    total_votes_01 AS "Suara Paslon 01",
    total_votes_02 AS "Suara Paslon 02",
    total_votes_03 AS "Suara Paslon 03",
    total_valid_votes AS "Jumlah Suara Sah",
    total_invalid_votes AS "Jumlah Suara Tidak Sah",
    total_votes AS "Jumlah Seluruh Suara",
    jml_hak_pilih AS "Jumlah Hak Pilih",
    
    link AS "Link", updated_at AS "Diperbarui Pada",
    obtained_at AS "Diperoleh Pada"
FROM %s WHERE (selisih_suara_paslon_dan_jumlah_sah <> 0 AND selisih_suara_sah_tidak_sah_dan_total <> 0) OR total_sum_votes > 300 OR total_votes > 300 OR jml_hak_pilih > 300 ORDER BY total_sum_votes+total_votes+jml_hak_pilih+selisih_suara_paslon_dan_jumlah_sah+selisih_suara_sah_tidak_sah_dan_total DESC LIMIT 5000`

const queryPotentialTableAllIn = `SELECT
	code AS "Kode",
	provinsi AS "Provinsi",
	kabupaten AS "Kabupaten",
	kecamatan AS "Kecamatan",
	kelurahan AS "Kelurahan",
	tps AS "TPS",
	total_votes_01 AS "Suara Paslon 01",
	total_votes_02 AS "Suara Paslon 02",
	total_votes_03 AS "Suara Paslon 03",
	
	total_valid_votes AS "Jumlah Suara Sah",
	total_invalid_votes AS "Jumlah Suara Tidak Sah",
	total_votes AS "Jumlah Seluruh Suara",
	jml_hak_pilih AS "Jumlah Hak Pilih",
	link AS "Link",
	updated_at AS "Diperbarui Pada",
	obtained_at AS "Diperoleh Pada"
FROM %s WHERE selisih_suara_paslon_dan_jumlah_sah = 0 AND selisih_suara_sah_tidak_sah_dan_total = 0 AND all_in > 0 AND total_sum_votes <= 300 AND total_votes <= 300 AND jml_hak_pilih <= 300 ORDER BY total_sum_votes DESC LIMIT 5000 `

func (r *repo) GetMapDist() ([]model.MapDist, error) {

	var mapDists = make([]model.MapDist, 0)
	row, err := r.db.Query(fmt.Sprintf(queryMapDist, r.tableRecord))
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		var mapDist model.MapDist

		if err := row.Scan(
			&mapDist.Winner,
			&mapDist.VotesValue,
			&mapDist.Chart.Paslon01,
			&mapDist.Chart.Paslon02,
			&mapDist.Chart.Paslon03,
			&mapDist.Latitude,
			&mapDist.Longitude,
			&mapDist.City,
		); err != nil {
			return nil, err
		}

		mapDists = append(mapDists, mapDist)
	}

	return mapDists, nil
}
func (r *repo) GetPotentialTableSum() ([]model.HHCWEntity, error) {
	var group = make([]model.HHCWEntity, 0)
	row, err := r.db.Query(fmt.Sprintf(queryPotentialTableSum, r.tableRecord))
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		var obj model.HHCWEntity
		var (
			obtainedAt int64
			updatedAt  int64
		)

		if err := row.Scan(
			&obj.Code,
			&obj.Provinsi,
			&obj.Kabupaten,
			&obj.Kecamatan,
			&obj.Kelurahan,
			&obj.TPS,
			&obj.Chart.Paslon01,
			&obj.Chart.Paslon02,
			&obj.Chart.Paslon03,
			&obj.Administrasi.Suara.Sah,
			&obj.Administrasi.Suara.TidakSah,
			&obj.Administrasi.Suara.Total,
			&obj.Administrasi.PenggunaTotal.Jumlah,
			&obj.Link,
			&updatedAt,
			&obtainedAt,
		); err != nil {
			return nil, err
		}

		obj.ObtainedAt = time.UnixMilli(obtainedAt)
		obj.UpdatedAt = time.UnixMilli(updatedAt)

		group = append(group, obj)
	}

	return group, nil
}
func (r *repo) GetPotentialTableAllIn() ([]model.HHCWEntity, error) {
	var group = make([]model.HHCWEntity, 0)
	row, err := r.db.Query(fmt.Sprintf(queryPotentialTableAllIn, r.tableRecord))
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		var obj model.HHCWEntity
		var (
			obtainedAt int64
			updatedAt  int64
		)

		if err := row.Scan(
			&obj.Code,
			&obj.Provinsi,
			&obj.Kabupaten,
			&obj.Kecamatan,
			&obj.Kelurahan,
			&obj.TPS,
			&obj.Chart.Paslon01,
			&obj.Chart.Paslon02,
			&obj.Chart.Paslon03,
			&obj.Administrasi.Suara.Sah,
			&obj.Administrasi.Suara.TidakSah,
			&obj.Administrasi.Suara.Total,
			&obj.Administrasi.PenggunaTotal.Jumlah,
			&obj.Link,
			&updatedAt,
			&obtainedAt,
		); err != nil {
			return nil, err
		}

		obj.ObtainedAt = time.UnixMilli(obtainedAt)
		obj.UpdatedAt = time.UnixMilli(updatedAt)

		group = append(group, obj)
	}

	return group, nil
}

var queryInsert = `
INSERT INTO %s (key, val)
VALUES %s
ON CONFLICT (key) DO UPDATE
SET val = EXCLUDED.val;
`

func (r *repo) UpdateStaticSummary(summaries []model.StaticSummary) error {
	var mapData = make([]interface{}, 0, len(summaries))
	for _, summary := range summaries {
		mapData = append(mapData, buildArgsStatic(summary)...)
	}
	q := fmt.Sprintf(queryInsert, r.tableKeyVal, buildPlaceholder(2, len(summaries)))
	_, err := r.db.Exec(q, mapData...)
	return err
}

func buildArgsStatic(summaries model.StaticSummary) []interface{} {
	return []interface{}{
		summaries.Key,
		summaries.Val,
	}
}
