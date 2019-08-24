package app

import (
	"bytes"
	"context"
	"io"

	"dflimg"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (a *App) GetS3File(ctx context.Context, resource *dflimg.Resource) ([]byte, error) {
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

	return bytes, nil
}
