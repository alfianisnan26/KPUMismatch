package dao

import "time"

type Cache interface {
	Get(key string, expiry time.Duration, fallback func(string) ([]byte, error)) ([]byte, error)
}
