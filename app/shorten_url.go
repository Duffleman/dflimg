package app

import (
	"context"
	"fmt"

	"dflimg"
	"dflimg/dflerr"
	"dflimg/rpc/middleware"

	"github.com/cuvva/ksuid"
)

// ShortenURL shortens a URL
func (a *App) ShortenURL(ctx context.Context, req *dflimg.CreateResourceRequest) (*dflimg.CreateResourceResponse, error) {
	username := ctx.Value(middleware.UsernameKey).(string)
	urlID := ksuid.Generate("url").String()

	err := a.db.FindShortcutConflicts(ctx, req.Shortcuts)
	if err != nil {
		return nil, dflerr.New("shortcuts already taken", dflerr.M{"shortcuts": req.Shortcuts}, dflerr.Parse(err))
	}

	// save to DB
	urlRes, err := a.db.NewURL(ctx, urlID, req.URL, username, req.Shortcuts, req.NSFW)
	if err != nil {
		return nil, err
	}

	rootURL := dflimg.GetEnv("root_url")
	hash := a.makeHash(urlRes.Serial)
	fullURL := fmt.Sprintf("%s/%s", rootURL, hash)

	return &dflimg.CreateResourceResponse{
		ResourceID: urlRes.ID,
		Type:       urlRes.Type,
		Hash:       hash,
		URL:        fullURL,
	}, nil
}
