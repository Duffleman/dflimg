package app

import (
	"bytes"
	"context"
	"dflimg"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-pg/pg"
)

// GetFile endpoint that gets a file by it's hash
func (a *App) GetFile(ctx context.Context, hash string) (*bytes.Buffer, error) {
	serial, err := a.decodeHash(hash)
	if err != nil {
		return nil, err
	}

	return a.getFileBySerial(ctx, serial)
}

func (a *App) getFileBySerial(ctx context.Context, serial int) (*bytes.Buffer, error) {
	file, err := a.db.FindFileBySerial(serial)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, dflimg.ErrNotFound
		}
		return nil, err
	}

	fileExt := getExtension(file.Type)
	fileKey := fmt.Sprintf("%s/%s%s", dflimg.S3RootKey, file.ID, fileExt)

	s3download, err := s3.New(a.aws).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(dflimg.S3Bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	io.Copy(&buf, s3download.Body)

	return &buf, nil
}

func (a *App) decodeHash(hash string) (int, error) {
	var set []int

	set, err := a.hasher.DecodeWithError(hash)
	if len(set) != 1 {
		return 0, errors.New("expecing 1 item in decoded hash")
	}

	return set[0], err
}
