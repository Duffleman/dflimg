package app

import (
	"bytes"
	"context"
	"net/http"

	"dflimg/rpc/middleware"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cuvva/ksuid"
	"github.com/kr/pretty"
)

// Upload is an app method that takes in a file and stores it
func (a *App) Upload(ctx context.Context, fileContent bytes.Buffer) (interface{}, error) {
	// get user
	username := ctx.Value(middleware.UsernameKey).(string)
	contentType := http.DetectContentType(fileContent.Bytes())
	fileID := ksuid.Generate("file").String()

	// upload to S3
	s3file, err := s3.New(a.aws).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(S3Bucket),
		Key:           aws.String(fileID),
		ACL:           aws.String("private"),
		Body:          bytes.NewReader(fileContent.Bytes()),
		ContentLength: aws.Int64(int64(fileContent.Len())),
		ContentType:   aws.String(contentType),
	})

	pretty.Println(s3file)

	// save to DB
	err = a.db.NewFile(fileID, username, contentType)
	if err != nil {
		return nil, err
	}

	// return hash of serial

	return nil, err
}
