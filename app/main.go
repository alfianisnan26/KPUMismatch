package main

import (
	"fmt"
	"kawalrealcount/internal/data/dao"
	"kawalrealcount/internal/data/model"
	"kawalrealcount/internal/pkg/httpclient/kpu"
	"kawalrealcount/internal/pkg/postgresql"
	"kawalrealcount/internal/pkg/redis"
	"kawalrealcount/internal/service/contributor"
	"kawalrealcount/internal/service/scrapper_v2"
	"kawalrealcount/internal/service/updater"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	jobControl = make(chan struct{}, 1)
)

const secret = "MUST MANUALLY SET ON EVERY BUILD"

func main() {

	redisHost := os.Getenv("REDIS_HOST")
	redisPasswd := os.Getenv("REDIS_PASSWORD")
	postgresTableRecord := os.Getenv("POSTGRES_TABLE")
	postgresTableStats := os.Getenv("POSTGRES_TABLE_STATS")
	postgresTableWebStats := os.Getenv("POSTGRES_TABLE_WEB_STATS")
	postgresTableHistogram := os.Getenv("POSTGRES_TABLE_HISTOGRAM")
	postgresTableKeyVal := os.Getenv("POSTGRES_TABLE_KEY_VAL")
	maximumRunningThread, _ := strconv.Atoi(os.Getenv("MAX_RUNNING_THREAD"))
	batchInsertLength, _ := strconv.Atoi(os.Getenv("BATCH_INSERT_LENGTH"))

	postgresUrl := os.Getenv("POSTGRES_URL")
	cooldownMinutes, _ := strconv.Atoi(os.Getenv("COOLDOWN_MINUTES"))
	cooldownMinutesUpdate, _ := strconv.Atoi(os.Getenv("COOLDOWN_MINUTES_UPDATE"))
	coolDownMinutesUpdateStatic, _ := strconv.Atoi(os.Getenv("COOLDOWN_MINUTES_UPDATE_STATIC"))
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
	cacheRepo, err = redis.New(redis.Param{
		Host:     redisHost,
		Password: redisPasswd,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

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

	kpuRepo := kpu.New(kpu.Param{
		CacheRepo: cacheRepo,
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

	updaterSvc := updater.New(updater.Param{
		UpdaterDatabaseRepo: psql,
	})

	quit := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(3)

	go func(wg *sync.WaitGroup, quit <-chan struct{}) {
		defer wg.Done()
		for {
			if err := updaterSvc.UpdateStats(); err != nil {
				fmt.Println(err)
				return
			}

			sleep := time.After(time.Minute * time.Duration(cooldownMinutesUpdate))
		LoopFor1:
			for {
				select {
				case <-quit:
					return
				case <-sleep:
					break LoopFor1
				}
			}
		}
	}(&wg, quit)

	go func(wg *sync.WaitGroup, quit <-chan struct{}) {
		defer wg.Done()
		for {
			if err := svcv2.ScrapAll(); err != nil {
				fmt.Println(err)
				return
			}

			if cooldownMinutes == 0 {
				return
			}

			fmt.Printf("Sleeping for %d minutes", cooldownMinutes)

			sleep := time.After(time.Minute * time.Duration(cooldownMinutes))

		LoopFor2:
			for {
				select {
				case <-quit:
					return
				case <-sleep:
					break LoopFor2
				}
			}
		}
	}(&wg, quit)

	go func(wg *sync.WaitGroup, quit <-chan struct{}) {
		defer wg.Done()
		for {
			if err := updaterSvc.UpdateStaticStats(); err != nil {
				fmt.Println(err)
				return
			}

			sleep := time.After(time.Minute * time.Duration(coolDownMinutesUpdateStatic))
		LoopFor3:
			for {
				select {
				case <-quit:
					return
				case <-sleep:
					break LoopFor3
				}
			}
		}
	}(&wg, quit)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a termination signal
	<-sigs

	close(quit)

	wg.Wait()

	return
}
