package app

import (
	"context"
	"fmt"
	"time"

	"dflimg"
	"dflimg/dflerr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cuvva/ksuid"
	pkgerr "github.com/pkg/errors"
)

// CreateSignedURL creates a file resource, but instead of accepting the file
// it generates a signed URL
func (a *App) CreateSignedURL(ctx context.Context, username string, req *dflimg.CreateSignedURLRequest) (*dflimg.CreateSignedURLResponse, error) {
	fileID := ksuid.Generate("file").String()
	fileKey := fmt.Sprintf("%s/%s", dflimg.S3RootKey, fileID)

	err := a.db.FindShortcutConflicts(ctx, req.Shortcuts)
	if err != nil {
		return nil, dflerr.New("shortcut_conflict", dflerr.M{"shortcuts": req.Shortcuts}, dflerr.Parse(err))
	}

	fileRes, err := a.db.NewPendingFile(ctx, fileID, fileKey, req.ContentType, username, req.Shortcuts, req.NSFW)
	if err != nil {
		return nil, err
	}

	s3req, _ := s3.New(a.aws).PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(dflimg.S3Bucket),
		Key:         aws.String(fileKey),
		ContentType: aws.String(req.ContentType),
	})

	url, err := s3req.Presign(15 * time.Minute)
	if err != nil {
		return nil, pkgerr.Wrap(err, "unable to create presigned s3 url")
	}

	rootURL := dflimg.GetEnv("root_url")
	hash := a.makeHash(fileRes.Serial)
	fullURL := fmt.Sprintf("%s/%s", rootURL, hash)

	gctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go a.saveHash(gctx, cancel, fileRes.Serial, hash)

	return &dflimg.CreateSignedURLResponse{
		ResourceID: fileRes.ID,
		Type:       fileRes.Type,
		Hash:       hash,
		URL:        fullURL,
		S3Link:     url,
	}, nil
}

func (a *App) saveHash(ctx context.Context, c context.CancelFunc, serial int, hash string) error {
	defer c()

	return a.db.SaveHash(ctx, serial, hash)
}
