package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"kawalrealcount/internal/data/model"
	"kawalrealcount/internal/pkg/httpclient/kpu"
	"kawalrealcount/internal/pkg/postgresql"
	"kawalrealcount/internal/service/contributor"
	"kawalrealcount/internal/service/scrapper_v2"
	"kawalrealcount/internal/service/updater"
)

const secret = "MUST MANUALLY SET ON EVERY BUILD"

const (
	workerTypeScrapper       = "SCRAPPER"
	workerTypeStatsUpdater   = "STATS_UPDATER"
	workerTypeStatisticStats = "STATISTIC_STATS_UPDATER"
	workerTypeTPSSync        = "TPS_SYNC"
)

func main() {

	postgresTableRecord := os.Getenv("POSTGRES_TABLE")
	postgresTableStats := os.Getenv("POSTGRES_TABLE_STATS")
	postgresTableWebStats := os.Getenv("POSTGRES_TABLE_WEB_STATS")
	postgresTableHistogram := os.Getenv("POSTGRES_TABLE_HISTOGRAM")
	mrt := os.Getenv("MAX_RUNNING_THREAD")
	postgresTableKeyVal := strings.Trim(os.Getenv("POSTGRES_TABLE_KEY_VAL"), `" \n`)
	maximumRunningThread, errMrt := strconv.Atoi(mrt)
	if errMrt != nil {
		panic(errMrt)
	}
	bil := strings.Trim(os.Getenv("BATCH_INSERT_LENGTH"), `" \n`)
	batchInsertLength, errBil := strconv.Atoi(bil)
	if errBil != nil {
		panic(errBil)
	}
	workerType := os.Getenv("WORKER_TYPE")
	postgresUrl := os.Getenv("POSTGRES_URL")
	contributionToken := os.Getenv("CONTRIBUTION_TOKEN")

	var contributorData model.ContributorData
	if contributionToken != "" {
		contributionSvc, err := contributor.New(contributor.Param{
			Secret: secret,
		})
		if err == nil {
			contributorData, err = contributionSvc.FetchContributionData(contributionToken)
			if err == nil {
				//redisHost = contributorData.RedisHost
				postgresTableRecord = contributorData.PostgresTableRecord
				postgresTableStats = contributorData.PostgresTableStats
				postgresUrl = contributorData.PostgresUrl
			} else {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	var (
		err error
	)

	psql, err := postgresql.New(postgresql.Param{
		ConnectionURL:  postgresUrl,
		TableRecord:    postgresTableRecord,
		TableStats:     postgresTableStats,
		TableWebStats:  postgresTableWebStats,
		TableHistogram: postgresTableHistogram,
		TableKeyVal:    postgresTableKeyVal,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	kpuRepo := kpu.New()

	svcv2 := scrapper_v2.New(scrapper_v2.Param{
		KPURepo:                    kpuRepo,
		DatabaseRepo:               psql,
		MaximumRunningThread:       maximumRunningThread,
		BatchInsertLength:          batchInsertLength,
		Contributor:                contributorData,
		ValidRecordExpiry:          time.Hour * 3,
		NotNullInvalidRecordExpiry: time.Minute * 30,
		NullRecordExpiry:           time.Minute * 10,
	})

	updaterSvc := updater.New(updater.Param{
		UpdaterDatabaseRepo: psql,
	})

	switch workerType {
	case workerTypeStatsUpdater:
		if err := updaterSvc.UpdateStats(); err != nil {
			fmt.Println(err)
			return
		}
	case workerTypeStatisticStats:
		if err := updaterSvc.UpdateStaticStats(); err != nil {
			fmt.Println(err)
			return
		}
	case workerTypeTPSSync:
		if err := svcv2.ScrapAllPPWT(); err != nil {
			fmt.Println(err)
			return
		}
	case workerTypeScrapper:
		fallthrough
	default:
		if err := svcv2.ScrapAll(); err != nil {
			fmt.Println(err)
			return
		}
	}
}
