package kawalpemilu

import "kawalrealcount/internal/data/dao"

type repo struct {
	cacheRepo dao.Cache
}

type Param struct {
	CacheRepo dao.Cache
}

func New(param Param) dao.KawalPemilu {
	return &repo{
		cacheRepo: param.CacheRepo,
	}
}
