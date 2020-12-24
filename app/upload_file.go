package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"dflimg"
	"dflimg/rpc/middleware"

	"github.com/cuvva/ksuid-go"
)

// UploadFile is an app method that takes in a file and stores it
func (a *App) UploadFile(ctx context.Context, req *dflimg.CreateFileRequest) (*dflimg.CreateResourceResponse, error) {
	// get user
	username := ctx.Value(middleware.UsernameKey).(string)
	bytes := req.File.Bytes()
	contentType := http.DetectContentType(bytes)
	fileID := ksuid.Generate("file").String()
	fileKey := a.fileProvider.GenerateKey(fileID)
	name := req.Name

	// upload to the file provider
	err := a.fileProvider.Upload(ctx, fileKey, contentType, req.File)
	if err != nil {
		return nil, err
	}

	// save to DB
	fileRes, err := a.db.NewFile(ctx, fileID, fileKey, username, name, contentType)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("file/%s", fileRes.Link)
	now := time.Now()

	if len(bytes) < MaxCacheSize {
		a.redis.Set(cacheKey, &CacheItem{
			Content: bytes,
			ModTime: &now,
		})
	}

	rootURL := dflimg.GetEnv("root_url")
	hash := a.makeHash(fileRes.Serial)
	fullURL := fmt.Sprintf("%s/%s", rootURL, hash)

	gctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go a.saveHash(gctx, cancel, fileRes.Serial, hash)

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
