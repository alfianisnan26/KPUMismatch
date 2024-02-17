package main

import (
	"fmt"
	"kawalrealcount/internal/data/dao"
	"kawalrealcount/internal/data/model"
	"kawalrealcount/internal/pkg/httpclient/kawalpemilu"
	"kawalrealcount/internal/pkg/httpclient/kpu"
	"kawalrealcount/internal/pkg/postgresql"
	"kawalrealcount/internal/pkg/redis"
	"kawalrealcount/internal/service/contributor"
	"kawalrealcount/internal/service/scrapper"
	"kawalrealcount/internal/service/scrapper_v2"
	"os"
	"strconv"
	"time"
)

var (
	jobControl = make(chan struct{}, 1)
)

const secret = "MUST MANUALLY SET ON EVERY BUILD"

func main() {

	filePath := os.Getenv("FILE_PATH")
	noCache := os.Getenv("NO_CACHE") == "True"
	redisHost := os.Getenv("REDIS_HOST")
	postgresTableRecord := os.Getenv("POSTGRES_TABLE")
	postgresTableStats := os.Getenv("POSTGRES_TABLE_STATS")
	postgresTableWebStats := os.Getenv("POSTGRES_TABLE_WEB_STATS")
	maximumRunningThread, _ := strconv.Atoi(os.Getenv("MAX_RUNNING_THREAD"))
	batchInsertLength, _ := strconv.Atoi(os.Getenv("BATCH_INSERT_LENGTH"))

	postgresUrl := os.Getenv("POSTGRES_URL")
	cooldownMinutes, _ := strconv.Atoi(os.Getenv("COOLDOWN_MINUTES"))
	scrapAll := os.Getenv("SCRAP_ALL") == "True"

	contributionToken := os.Getenv("CONTRIBUTION_TOKEN")

	var contributorData model.ContributorData
	if contributionToken != "" {
		contributionSvc, err := contributor.New(contributor.Param{
			Secret: secret,
		})
		if err == nil {
			contributorData, err = contributionSvc.FetchContributionData(contributionToken)
			if err == nil {
				redisHost = contributorData.RedisHost
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
		cacheRepo dao.Cache
		err       error
	)
	if !noCache {
		cacheRepo, err = redis.New(redis.Param{
			Host: redisHost,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	psql, err := postgresql.New(postgresql.Param{
		ConnectionURL: postgresUrl,
		TableRecord:   postgresTableRecord,
		TableStats:    postgresTableStats,
		TableWebStats: postgresTableWebStats,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	kpuRepo := kpu.New(kpu.Param{
		CacheRepo: cacheRepo,
	})
	kawalPemiluRepo := kawalpemilu.New(kawalpemilu.Param{
		CacheRepo: cacheRepo,
	})

	svc := scrapper.New(scrapper.Param{
		Contributor:     contributorData,
		KPURepo:         kpuRepo,
		KawalPemiluRepo: kawalPemiluRepo,
		DatabaseRepo:    psql,

		MaximumRunningThread: maximumRunningThread,
		BatchInsertLength:    batchInsertLength,
	})

	svcv2 := scrapper_v2.New(scrapper_v2.Param{
		KPURepo:                    kpuRepo,
		CacheRepo:                  cacheRepo,
		DatabaseRepo:               psql,
		MaximumRunningThread:       maximumRunningThread,
		BatchInsertLength:          batchInsertLength,
		Contributor:                contributorData,
		ValidRecordExpiry:          time.Hour * 3,
		NotNullInvalidRecordExpiry: time.Minute * 30,
		NullRecordExpiry:           time.Minute * 10,
	})

	// first run
	var fn func()

	if scrapAll {
		fn = func() {
			err := svcv2.ScrapAll()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	} else {
		fn = func() {
			if err := svc.ScrapAllCompiled(filePath, false); err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}

	for {

		fn()

		if cooldownMinutes == 0 {
			return
		}

		fmt.Printf("Sleeping for %d minutes", cooldownMinutes)

		time.Sleep(time.Minute * time.Duration(cooldownMinutes))
	}

}
