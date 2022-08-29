package cache

import (
	"time"
)

const (
	cacheBaseDir = "/tmp/sort-awesome-go-cache"
	ttl          = time.Hour * 24
)

type Handler interface {
	ReadString(key string) string
	Has(key string) bool
	EraseAll() error
	Write(key string, val []byte) error
}
