package app

import (
	"context"
	"dflimg/dflerr"
	"fmt"
	"time"

	"dflimg"

	"github.com/cuvva/ksuid-go"
	pkgerr "github.com/pkg/errors"
)

// CreateSignedURL creates a file resource, but instead of accepting the file
// it generates a signed URL
func (a *App) CreateSignedURL(ctx context.Context, username string, contentType string) (*dflimg.CreateSignedURLResponse, error) {
	if !a.fileProvider.SupportsTwoStage() {
		return nil, dflerr.New("signed_urls_unsupported", nil)
	}

	fileID := ksuid.Generate("file").String()
	fileKey := a.fileProvider.GenerateKey(fileID)

	fileRes, err := a.db.NewPendingFile(ctx, fileID, fileKey, username, contentType)
	if err != nil {
		return nil, err
	}

	url, err := a.fileProvider.PrepareUpload(ctx, fileKey, contentType, 15*time.Minute)
	if err != nil {
		return nil, pkgerr.Wrap(err, "unable to create presigned url")
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
