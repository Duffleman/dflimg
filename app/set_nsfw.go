package app

import (
	"context"
)

func (a *App) SetNSFW(ctx context.Context, resourceID string, nsfw bool) error {
	return a.db.SetNSFW(ctx, resourceID, nsfw)
}
