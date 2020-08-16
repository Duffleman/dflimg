package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dflimg"
	"dflimg/dflerr"
)

const maxSizeInMB = 32
const byteJump = 1024

// MaxCacheSize is the maximum size of a file for it to skip the cache: 536,870,912
const MaxCacheSize = byteJump * byteJump * maxSizeInMB

// GetFile returns a file from the cache,or the file provider
func (a *App) GetFile(ctx context.Context, resource *dflimg.Resource) ([]byte, *time.Time, error) {
	cacheKey := fmt.Sprintf("file/%s", resource.Link)

	if item, found := a.redis.Get(cacheKey); found {
		return item.Content, item.ModTime, nil
	}

	bytes, lastModified, err := a.fileProvider.Get(ctx, resource)
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchKey") {
			return nil, nil, dflerr.ErrNotFound
		}

		return nil, nil, err
	}

	if len(bytes) < MaxCacheSize {
		a.redis.Set(cacheKey, &CacheItem{
			Content: bytes,
			ModTime: lastModified,
		})
	}

	return bytes, lastModified, nil
}
