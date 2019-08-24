package app

import (
	"context"
	"fmt"

	"dflimg"
	"dflimg/dflerr"
	"dflimg/rpc/middleware"

	"github.com/cuvva/ksuid"
)

func (a *App) ShortenURL(ctx context.Context, url string, shortcuts []string) (*dflimg.ResponseCreatedResponse, error) {
	username := ctx.Value(middleware.UsernameKey).(string)
	urlID := ksuid.Generate("url").String()

	err := a.db.FindShortcutConflicts(ctx, shortcuts)
	if err != nil {
		return nil, dflerr.New("shortcuts already taken", dflerr.M{"shortcuts": shortcuts}, dflerr.Parse(err))
	}

	// save to DB
	urlRes, err := a.db.NewURL(ctx, urlID, url, username, shortcuts)
	if err != nil {
		return nil, err
	}

	rootURL := dflimg.GetEnv("root_url")
	hash := a.makeHash(urlRes.Serial)
	fullURL := fmt.Sprintf("%s/%s", rootURL, hash)

	return &dflimg.ResponseCreatedResponse{
		ResourceID: urlRes.ID,
		Type:       urlRes.Type,
		Hash:       hash,
		URL:        fullURL,
	}, nil
}
