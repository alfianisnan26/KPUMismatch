package main

import (
	"flag"
	"fmt"
	"kawalrealcount/internal/data/dao"
	"kawalrealcount/internal/pkg/httpclient/kawalpemilu"
	"kawalrealcount/internal/pkg/httpclient/kpu"
	"kawalrealcount/internal/pkg/redis"
	"kawalrealcount/internal/pkg/sqlite"
	"kawalrealcount/internal/service/scrapper"
)

func main() {

	filePath := flag.String("filepath", "report.xlsx", "set output filepath")
	noCache := flag.Bool("nocache", false, "use only sql db")
	redisHost := flag.String("redishost", "localhost:6379", "set redis host")
	sqlitePath := flag.String("sqlitepath", "db.sqlite3", "set sqlite path")
	var (
		cacheRepo dao.Cache
		err       error
	)
	if !*noCache {
		cacheRepo, err = redis.New(redis.Param{
			Host: *redisHost,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	if cacheRepo == nil {
		cacheRepo, err = sqlite.New(sqlite.Param{
			FilePath: *sqlitePath,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	kpuRepo := kpu.New(kpu.Param{
		CacheRepo: cacheRepo,
	})
	kawalPemiluRepo := kawalpemilu.New(kawalpemilu.Param{
		CacheRepo: cacheRepo,
	})

	svc := scrapper.New(scrapper.Param{
		KPURepo:              kpuRepo,
		KawalPemiluRepo:      kawalPemiluRepo,
		MaximumRunningThread: 50,
	})

	if err := svc.ScrapAllCompiled(*filePath); err != nil {
		fmt.Println(err.Error())
		return
	}

}
