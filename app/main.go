package main

import (
	"fmt"
	"github.com/robfig/cron"
	"kawalrealcount/internal/data/dao"
	"kawalrealcount/internal/pkg/httpclient/kawalpemilu"
	"kawalrealcount/internal/pkg/httpclient/kpu"
	"kawalrealcount/internal/pkg/postgresql"
	"kawalrealcount/internal/pkg/redis"
	"kawalrealcount/internal/pkg/sqlite"
	"kawalrealcount/internal/service/scrapper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	filePath := os.Getenv("FILE_PATH")
	noCache := os.Getenv("NO_CACHE") == "True"
	redisHost := os.Getenv("REDIS_HOST")
	sqlitePath := os.Getenv("SQLITE_PATH")
	postgresTable := os.Getenv("POSTGRES_TABLE")
	postgresUrl := os.Getenv("POSTGRES_URL")
	schedulePattern := os.Getenv("SCHEDULE_PATTERN")

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
		TableName:     postgresTable,
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
		KPURepo:         kpuRepo,
		KawalPemiluRepo: kawalPemiluRepo,
		DatabaseRepo:    psql,

		MaximumRunningThread: 50,
	})

	// first run
	fn := func() {
		if err := svc.ScrapAllCompiled(filePath); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	fn()

	cronJob := cron.NewWithLocation(time.Local)

	if err := cronJob.AddFunc(schedulePattern, fn); err != nil {
		fmt.Println(err.Error())
		return
	}

	cronJob.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	cronJob.Stop()

	os.Exit(0)
}
