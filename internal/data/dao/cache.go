package dao

import (
	"kawalrealcount/internal/data/model"
	"time"
)

type Cache interface {
	Get(key string, expiry time.Duration, fallback func(string) ([]byte, error)) ([]byte, error)
	GetHHCW(key string) (model.HHCWEntity, error)
	PutHHCW(key string, data model.HHCWEntity, expiry time.Duration) error
}
