package main

import (
	"fmt"
	"kawalrealcount/internal/data/dao"
	"kawalrealcount/internal/data/model"
	"kawalrealcount/internal/pkg/httpclient/kawalpemilu"
	"kawalrealcount/internal/pkg/httpclient/kpu"
	"kawalrealcount/internal/pkg/postgresql"
	"kawalrealcount/internal/pkg/redis"
	"kawalrealcount/internal/pkg/sqlite"
	"kawalrealcount/internal/service/contributor"
	"kawalrealcount/internal/service/scrapper"
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
	sqlitePath := os.Getenv("SQLITE_PATH")
	postgresTableRecord := os.Getenv("POSTGRES_TABLE")
	postgresTableStats := os.Getenv("POSTGRES_TABLE_STATS")
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

	if cacheRepo == nil {
		cacheRepo, err = sqlite.New(sqlite.Param{
			FilePath: sqlitePath,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	psql, err := postgresql.New(postgresql.Param{
		ConnectionURL: postgresUrl,
		TableRecord:   postgresTableRecord,
		TableStats:    postgresTableStats,
	})
	if err != nil {
		fmt.Println(err.Error())
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

		MaximumRunningThread: 15,
		ProgressRefreshRate:  time.Minute,
	})

	// first run
	var fn func()

	if scrapAll {
		fn = func() {
			if err := svc.ScrapAllSeedOnly(); err != nil {
				fmt.Println(err.Error())
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
