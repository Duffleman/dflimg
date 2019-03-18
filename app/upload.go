package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"dflimg"
	"dflimg/rpc/middleware"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cuvva/ksuid"
)

// Upload is an app method that takes in a file and stores it
func (a *App) Upload(ctx context.Context, fileContent bytes.Buffer, labels []string) (*dflimg.UploadFileResponse, error) {
	// get user
	username := ctx.Value(middleware.UsernameKey).(string)
	contentType := http.DetectContentType(fileContent.Bytes())
	fileID := ksuid.Generate("file").String()
	fileExt := getExtension(contentType)
	fileKey := fmt.Sprintf("%s/%s%s", dflimg.S3RootKey, fileID, fileExt)

	for _, label := range labels {
		// TODO: make more efficient
		_, err := a.db.FindFileByLabel(label)
		if err == nil {
			return nil, errors.New("label already taken")
		}
	}

	// upload to S3
	_, err := s3.New(a.aws).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(dflimg.S3Bucket),
		Key:           aws.String(fileKey),
		ACL:           aws.String("private"),
		Body:          bytes.NewReader(fileContent.Bytes()),
		ContentLength: aws.Int64(int64(fileContent.Len())),
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		return nil, err
	}

	// save to DB
	err = a.db.NewFile(fileID, fileKey, username, contentType, labels)
	if err != nil {
		return nil, err
	}

	file, err := a.db.FindFile(fileID)
	if err != nil {
		return nil, err
	}

	rootURL := dflimg.GetEnv("root_url")
	hash := a.makeHash(file.Serial)
	url := fmt.Sprintf("%s/%s", rootURL, hash)

	return &dflimg.UploadFileResponse{
		FileID: file.ID,
		Hash:   hash,
		URL:    url,
	}, err
}

func getExtension(t string) string {
	switch t {
	case "text/plain":
		return ".txt"
	default:
		ext := strings.Split(t, "/")
		return fmt.Sprintf(".%s", ext[1])
	}
}

func (a *App) makeHash(serial int) string {
	e, _ := a.hasher.Encode([]int{serial})

	return e
}
