# Kawal Real Count Pemilu Indonesia 2024

This worker is designed to fetch mismatch data from the official website of the Indonesian General Election Commission (KPU) for the 2024 elections, accessible via the open Sirekap HTTP API.

## Description

The purpose of this Golang worker is to retrieve data from the Sirekap API and store it in a PostgreSQL database. The worker performs the following tasks:

- Reconciling the sum of all scanned input total votes for three candidate values to the scanned input of legitimate votes, ensuring there is no difference.
- Additional data collection for All-In data (unique data) for total votes only intended for one candidate pair at one TPS (voting booth).

## How to Run

```sh
docker pull alfianisnan26/kawalrealcount:latest
docker run -d --env-file .env --name kawalrealcount alfianisnan26/kawalrealcount:latest
```

### Environment Variables

```env
FILE_PATH=report.xlsx
NO_CACHE=True
REDIS_HOST=localhost:6379
SQLITE_PATH=db.sqlite3
POSTGRES_TABLE=kpu_tps
POSTGRES_TABLE_STATS=kpu_tps_stats
POSTGRES_URL=postgres://admin:root@localhost:5432/postgres
SCHEDULE_PATTERN=0 */3 * * *
SCRAP_ALL=False
```

- `FILE_PATH`: Path to the Excel report file.
- `NO_CACHE`: Set to `True` to disable caching.
- `REDIS_HOST`: Redis server host address.
- `SQLITE_PATH`: Path to the SQLite database file.
- `POSTGRES_TABLE`: Name of the PostgreSQL table to store data.
- `POSTGRES_TABLE_STATS`: Name of the PostgreSQL tabe to store worker stats data.
- `POSTGRES_URL`: URL for connecting to the PostgreSQL database.
- `SCHEDULE_PATTERN`: Cron-like schedule pattern for periodic execution.
- `SCRAP_ALL`: To enable scrap all data without check the reconciliation process

## Disclaimer

This worker is provided for informational purposes only. Please handle data in accordance with relevant laws, regulations, and policies. Use it at your own discretion and risk.

## Credit & Documentation

- https://hub.docker.com/repository/docker/alfianisnan26/kawalrealcount
- https://github.com/alfianisnan26/KPUMismatch
- https://kawalsuara.alvilab.my.id by https://github.com/alvimuh