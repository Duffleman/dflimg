package app

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"dflimg"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	cache "github.com/patrickmn/go-cache"
)

// GetS3File returns a file from the cache, or S3
func (a *App) GetS3File(ctx context.Context, resource *dflimg.Resource) ([]byte, error) {
	cacheKey := fmt.Sprintf("file/%s", resource.Link)

	if file, found := a.cache.Get(cacheKey); found {
		return file.([]byte), nil
	}

	s3download, err := s3.New(a.aws).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(dflimg.S3Bucket),
		Key:    aws.String(resource.Link),
	})
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	_, err = io.Copy(&buf, s3download.Body)
	if err != nil {
		return nil, err
	}

	bytes := buf.Bytes()

	a.cache.Set(cacheKey, bytes, cache.DefaultExpiration)

	return bytes, nil
}
