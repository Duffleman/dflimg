package app

import (
	"context"
	"fmt"
	"time"

	"dflimg"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cuvva/ksuid"
	pkgerr "github.com/pkg/errors"
)

// CreatedSignedURL creates a file resource, but instead of accepting the file
// it generates a signed URL
func (a *App) CreatedSignedURL(ctx context.Context, username, contentType string, contentLength int, shortcuts []string, nsfw bool) (*dflimg.CreateSignedURLResponse, error) {
	fileID := ksuid.Generate("file").String()
	fileKey := fmt.Sprintf("%s/%s", dflimg.S3RootKey, fileID)

	_, err := a.db.NewPendingFile(ctx, fileID, fileKey, contentType, username, shortcuts, nsfw)
	if err != nil {
		return nil, err
	}

	req, _ := s3.New(a.aws).PutObjectRequest(&s3.PutObjectInput{
		Bucket:        aws.String(dflimg.S3Bucket),
		Key:           aws.String(fileKey),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(int64(contentLength)),
	})

	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		return nil, pkgerr.Wrap(err, "unable to create presigned s3 url")
	}

	return &dflimg.CreateSignedURLResponse{
		URL: url,
	}, nil
}
