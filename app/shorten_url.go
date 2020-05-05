package app

import (
	"context"
	"fmt"
	"time"

	"dflimg"
	"dflimg/rpc/middleware"

	"github.com/cuvva/ksuid-go"
)

// ShortenURL shortens a URL
func (a *App) ShortenURL(ctx context.Context, url string) (*dflimg.CreateResourceResponse, error) {
	username := ctx.Value(middleware.UsernameKey).(string)
	urlID := ksuid.Generate("url").String()

	// save to DB
	urlRes, err := a.db.NewURL(ctx, urlID, username, url)
	if err != nil {
		return nil, err
	}

	rootURL := dflimg.GetEnv("root_url")
	hash := a.makeHash(urlRes.Serial)
	fullURL := fmt.Sprintf("%s/%s", rootURL, hash)

	gctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go a.saveHash(gctx, cancel, urlRes.Serial, hash)

	return &dflimg.CreateResourceResponse{
		ResourceID: urlRes.ID,
		Type:       urlRes.Type,
		Hash:       hash,
		URL:        fullURL,
	}, nil
}
