--
-- Data Terekam					count(*)
-- Data Tidak Kosong (TPS, SUARA) 	count(*),sum(total_sum_votes),sum(total_votes),sum(jml_hak_pilih) when total_sum_votes > 0 OR total_votes > 0 OR jml_hak_pilih > 0
-- Data Kosong						Data Terekam - Data Tidak Kosong
-- Data Tidak Valid (TPS)				count(*),sum() when total_sum_votes != total_valid_votes OR total_valid_votes + total_invalid_votes != total_votes OR jml_hak_pilih > 300 OR total_votes > 300 OR total_sum_votes > 300
-- Data Valid (TPS)					Data Tidak Kosong - Data Valid

SELECT json_agg(t), json_agg()
FROM(
        SELECT
            count(*) as all_data,

            sum(CASE WHEN total_sum_votes > 0 OR total_votes > 0 OR jml_hak_pilih > 0 THEN 1
                     ELSE 0 END) AS not_null_count,

            sum(CASE WHEN total_sum_votes > 0 OR total_votes > 0 OR jml_hak_pilih > 0 THEN total_sum_votes + total_votes + jml_hak_pilih
                     ELSE 0 END) AS not_null_sum,

            sum(CASE WHEN total_sum_votes != total_valid_votes OR total_valid_votes + total_invalid_votes != total_votes OR total_votes > 300 OR jml_hak_pilih > 300 OR total_sum_votes > 300 THEN 1
                     ELSE 0 END) as not_valid_count,

            sum(CASE WHEN total_sum_votes != total_valid_votes OR total_valid_votes + total_invalid_votes != total_votes OR total_votes > 300 OR jml_hak_pilih > 300 OR total_sum_votes > 300
                         THEN GREATEST(ABS(total_sum_votes - total_valid_votes), ABS((total_valid_votes + total_invalid_votes) - total_votes), total_votes - 300, jml_hak_pilih - 300, total_sum_votes - 300)
                     ELSE 0 END) as not_valid_diff,

            sum(total_votes_01) as sum_raw_votes_01,

            sum(CASE WHEN NOT(total_sum_votes != total_valid_votes OR total_valid_votes + total_invalid_votes != total_votes OR total_votes > 300 OR jml_hak_pilih > 300 OR total_sum_votes > 300)
                         THEN total_votes_01 ELSE 0 END
            ) AS sum_votes_01,

            sum(CASE WHEN NOT(total_sum_votes != total_valid_votes OR total_valid_votes + total_invalid_votes != total_votes OR total_votes > 300 OR jml_hak_pilih > 300 OR total_sum_votes > 300)
                         THEN (CASE WHEN GREATEST(total_votes_01, total_votes_01, total_votes_03) = total_votes_01 THEN 1 ELSE 0 END) ELSE 0 END
            ) AS count_votes_01,

            sum(CASE WHEN NOT(total_sum_votes != total_valid_votes OR total_valid_votes + total_invalid_votes != total_votes OR total_votes > 300 OR jml_hak_pilih > 300 OR total_sum_votes > 300)
                         THEN (CASE WHEN total_sum_votes = total_votes_01 THEN 1 ELSE 0 END) ELSE 0 END
            ) AS count_all_in_01,

            sum(total_votes_02) as sum_raw_votes_02,

            sum(CASE WHEN NOT(total_sum_votes != total_valid_votes OR total_valid_votes + total_invalid_votes != total_votes OR total_votes > 300 OR jml_hak_pilih > 300 OR total_sum_votes > 300)
                         THEN total_votes_02 ELSE 0 END
            ) AS sum_votes_02,

            sum(CASE WHEN NOT(total_sum_votes != total_valid_votes OR total_valid_votes + total_invalid_votes != total_votes OR total_votes > 300 OR jml_hak_pilih > 300 OR total_sum_votes > 300)
                         THEN (CASE WHEN GREATEST(total_votes_01, total_votes_01, total_votes_03) = total_votes_02 THEN 1 ELSE 0 END) ELSE 0 END
            ) AS count_votes_02,

            sum(CASE WHEN NOT(total_sum_votes != total_valid_votes OR total_valid_votes + total_invalid_votes != total_votes OR total_votes > 300 OR jml_hak_pilih > 300 OR total_sum_votes > 300)
                         THEN (CASE WHEN total_sum_votes = total_votes_02 THEN 1 ELSE 0 END) ELSE 0 END
            ) AS count_all_in_02,

            sum(total_votes_03) as sum_raw_votes_03,

            sum(CASE WHEN NOT(total_sum_votes != total_valid_votes OR total_valid_votes + total_invalid_votes != total_votes OR total_votes > 300 OR jml_hak_pilih > 300 OR total_sum_votes > 300)
                         THEN total_votes_01 ELSE 0 END
            ) AS sum_votes_03,

            sum(CASE WHEN NOT(total_sum_votes != total_valid_votes OR total_valid_votes + total_invalid_votes != total_votes OR total_votes > 300 OR jml_hak_pilih > 300 OR total_sum_votes > 300)
                         THEN (CASE WHEN GREATEST(total_votes_01, total_votes_01, total_votes_03) = total_votes_03 THEN 1 ELSE 0 END) ELSE 0 END
            ) AS count_votes_03,

            sum(CASE WHEN NOT(total_sum_votes != total_valid_votes OR total_valid_votes + total_invalid_votes != total_votes OR total_votes > 300 OR jml_hak_pilih > 300 OR total_sum_votes > 300)
                         THEN (CASE WHEN total_sum_votes = total_votes_03 THEN 1 ELSE 0 END) ELSE 0 END
            ) AS count_all_in_03,

            SUBSTRING(code, 0, 5) as city_id

        FROM kpu_tps_v2 GROUP BY city_id
    ) as t