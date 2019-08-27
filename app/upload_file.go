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
	cache "github.com/patrickmn/go-cache"
)

// UploadFile is an app method that takes in a file and stores it
func (a *App) UploadFile(ctx context.Context, req *dflimg.CreateResourceRequest) (*dflimg.CreateResourceResponse, error) {
	// get user
	username := ctx.Value(middleware.UsernameKey).(string)
	contentType := http.DetectContentType(req.File.Bytes())
	fileID := ksuid.Generate("file").String()
	fileKey := fmt.Sprintf("%s/%s", dflimg.S3RootKey, fileID)

	err := a.db.FindShortcutConflicts(ctx, req.Shortcuts)
	if err != nil {
		return nil, dflerr.New("shortcuts already taken", dflerr.M{"shortcuts": req.Shortcuts}, dflerr.Parse(err))
	}

	// upload to S3
	_, err = s3.New(a.aws).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(dflimg.S3Bucket),
		Key:           aws.String(fileKey),
		ACL:           aws.String("private"),
		Body:          bytes.NewReader(req.File.Bytes()),
		ContentLength: aws.Int64(int64(req.File.Len())),
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		return nil, err
	}

	// save to DB
	fileRes, err := a.db.NewFile(ctx, fileID, fileKey, contentType, username, req.Shortcuts, req.NSFW)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("file/%s", fileRes.Link)
	a.cache.Set(cacheKey, req.File.Bytes(), cache.DefaultExpiration)

	rootURL := dflimg.GetEnv("root_url")
	hash := a.makeHash(fileRes.Serial)
	fullURL := fmt.Sprintf("%s/%s", rootURL, hash)

	return &dflimg.CreateResourceResponse{
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
