package app

import (
	"context"

	"dflimg"
	"dflimg/db"
)

func (a *App) RemoveShortcut(ctx context.Context, resource *dflimg.Resource, shortcut string) error {
	return a.db.ChangeShortcut(ctx, db.ArrayRemove, resource.ID, shortcut)
}
