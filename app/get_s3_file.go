package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"dflimg"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	cache "github.com/patrickmn/go-cache"
)

type CacheItem struct {
	Content []byte
	ModTime *time.Time
}

// GetS3File returns a file from the cache, or S3
func (a *App) GetS3File(ctx context.Context, resource *dflimg.Resource) ([]byte, *time.Time, error) {
	cacheKey := fmt.Sprintf("file/%s", resource.Link)

	if item, found := a.cache.Get(cacheKey); found {
		i := item.(*CacheItem)

		return i.Content, i.ModTime, nil
	}

	s3item, err := s3.New(a.aws).GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(dflimg.S3Bucket),
		Key:    aws.String(resource.Link),
	})
	if err != nil {
		return nil, nil, err
	}

	var buf bytes.Buffer

	_, err = io.Copy(&buf, s3item.Body)
	if err != nil {
		return nil, nil, err
	}

	bytes := buf.Bytes()

	a.cache.Set(cacheKey, &CacheItem{
		Content: bytes,
		ModTime: s3item.LastModified,
	}, cache.DefaultExpiration)

	return bytes, s3item.LastModified, nil
}
