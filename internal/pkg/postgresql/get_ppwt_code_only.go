package postgresql

import (
	"fmt"
	"kawalrealcount/internal/data/model"
	"time"
)

func (r *repo) GetPPWTCodeOnly(ppwtCh chan<- *model.HHCWEntity) error {

	res, err := r.db.Query(fmt.Sprintf(`SELECT id, code AS "Kode",
    total_votes_01 AS "Suara Paslon 01",
    total_votes_02 AS "Suara Paslon 02",
    total_votes_03 AS "Suara Paslon 03",
    total_valid_votes AS "Jumlah Suara Sah",
    total_invalid_votes AS "Jumlah Suara Tidak Sah",
    total_votes AS "Jumlah Seluruh Suara",
    jml_hak_pilih AS "Jumlah Hak Pilih",
    updated_at AS "Diperbarui Pada",
    obtained_at AS "Diperoleh Pada"
FROM %s`, r.tableRecord))
	if err != nil {
		return err
	}

	for res.Next() {
		var (
			obj                   model.HHCWEntity
			updatedAt, obtainedAt int64
			id                    int64
		)

		if err := res.Scan(
			&obj.ID,
			&obj.Code,
			&obj.Chart.Paslon01,
			&obj.Chart.Paslon02,
			&obj.Chart.Paslon03,
			&obj.Administrasi.Suara.Sah,
			&obj.Administrasi.Suara.TidakSah,
			&obj.Administrasi.Suara.Total,
			&obj.Administrasi.PenggunaTotal.Jumlah,
			&updatedAt,
			&obtainedAt,
		); err != nil {
			return err
		}

		obj.UpdatedAt = time.UnixMilli(updatedAt)
		obj.ObtainedAt = time.UnixMilli(obtainedAt)

		if len(obj.Code) < 13 {
			fmt.Println("Invalid Code", obj.Code, id)
			continue
		}

		ppwtCh <- &obj
	}
	return nil
}
