package main

import (
	"fmt"
	"kawalrealcount/internal/data/model"
	"kawalrealcount/internal/pkg/httpclient/kawalpemilu"
	"kawalrealcount/internal/pkg/httpclient/kpu"
	"kawalrealcount/internal/pkg/redis"
	"kawalrealcount/internal/service/scrapper"
)

func main() {

	cacheRepo, err := redis.New(redis.Param{
		Host: "localhost:6379",
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
		KPURepo:              kpuRepo,
		KawalPemiluRepo:      kawalPemiluRepo,
		MaximumRunningThread: 50,
	})

	criterion := model.Criterion{
		//IgnoreAll: true,

		//WithNonZeroMismatchSuara: true,
		//WithInvalidSumOfSuara:    true,

		//WithNonZeroMismatchSuaraAndPengguna: true,
		//WithDeltaErrThreshold: true,
		//DeltaErrThreshold:     50,

		//WithSumThreshold: true,
		//SumThreshold:     300,
		WithAllIn: true,
	}

	filePath := "export/all_in.csv"

	if err := svc.ScrapAllCompiled(criterion, filePath); err != nil {
		fmt.Println(err.Error())
		return
	}

}
