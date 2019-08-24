package app

import (
	"context"

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
	err = a.db.NewURL(ctx, urlID, url, username, shortcuts)
	if err != nil {
		return nil, err
	}

	return a.CreateResource(ctx, urlID)
}
