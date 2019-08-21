package app

import (
	"bytes"
	"context"
	"io"

	"dflimg"
	"dflimg/dflerr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-pg/pg"
)

// GetFileByHash gets a file by it's hash
func (a *App) GetFileByHash(ctx context.Context, hash string) (string, *bytes.Buffer, error) {
	serial, err := a.decodeHash(hash)
	if err != nil {
		return "", nil, err
	}

	return a.getFileBySerial(ctx, serial)
}

func (a *App) getFileBySerial(ctx context.Context, serial int) (string, *bytes.Buffer, error) {
	file, err := a.db.FindFileBySerial(ctx, serial)
	if err != nil {
		if err == pg.ErrNoRows {
			return "", nil, dflerr.New(dflerr.NotFound, nil)
		}
		return "", nil, err
	}

	s3download, err := s3.New(a.aws).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(dflimg.S3Bucket),
		Key:    aws.String(file.S3),
	})
	if err != nil {
		return "", nil, err
	}

	var buf bytes.Buffer
	io.Copy(&buf, s3download.Body)

	return file.Type, &buf, nil
}

func (a *App) decodeHash(hash string) (int, error) {
	var set []int

	set, err := a.hasher.DecodeWithError(hash)
	if len(set) != 1 {
		return 0, dflerr.New("cannot decode hash", dflerr.M{"hash": hash}, dflerr.New("expecting single hashed item in body", nil))
	}

	return set[0], err
}
