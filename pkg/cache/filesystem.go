package cache

import (
	"github.com/adelowo/onecache/filesystem"
	"os"
)

type FilesystemCache struct {
	store *filesystem.FSStore
}

func (c *FilesystemCache) ReadString(key string) string {
	val, err := c.store.Get(key)
	if err != nil {
		return ""
	}
	return string(val)
}

func (c *FilesystemCache) Has(key string) bool {
	return c.store.Has(key)
}

func (c *FilesystemCache) EraseAll() error {
	return c.store.Flush()
}

func (c *FilesystemCache) Write(key string, val []byte) error {
	return c.store.Set(key, val, ttl)
}

func NewLocalCache() (Handler, error) {
	err := os.MkdirAll(cacheBaseDir, 0744)
	if err != nil {
		return nil, err
	}
	return &FilesystemCache{store: filesystem.MustNewFSStore(cacheBaseDir)}, nil
}
