package app

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"dflimg"
	"dflimg/dflerr"
	"dflimg/rpc/middleware"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cuvva/ksuid"
)

// UploadFile is an app method that takes in a file and stores it
func (a *App) UploadFile(ctx context.Context, fileContent bytes.Buffer, shortcuts []string) (*dflimg.ResponseCreatedResponse, error) {
	// get user
	username := ctx.Value(middleware.UsernameKey).(string)
	contentType := http.DetectContentType(fileContent.Bytes())
	fileID := ksuid.Generate("file").String()
	fileKey := fmt.Sprintf("%s/%s", dflimg.S3RootKey, fileID)

	err := a.db.FindShortcutConflicts(ctx, shortcuts)
	if err != nil {
		return nil, dflerr.New("shortcuts already taken", dflerr.M{"shortcuts": shortcuts}, dflerr.Parse(err))
	}

	// upload to S3
	_, err = s3.New(a.aws).PutObject(&s3.PutObjectInput{
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
	fileRes, err := a.db.NewFile(ctx, fileID, fileKey, contentType, username, shortcuts)
	if err != nil {
		return nil, err
	}

	rootURL := dflimg.GetEnv("root_url")
	hash := a.makeHash(fileRes.Serial)
	fullURL := fmt.Sprintf("%s/%s", rootURL, hash)

	return &dflimg.ResponseCreatedResponse{
		ResourceID: fileRes.ID,
		Type:       fileRes.Type,
		Hash:       hash,
		URL:        fullURL,
	}, nil
}

func (a *App) makeHash(serial int) string {
	e, _ := a.hasher.Encode([]int{serial})

	return e
}
