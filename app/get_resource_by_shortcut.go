package app

import (
	"context"

	"dflimg"
)

func (a *App) GetResourceByShortcut(ctx context.Context, shortcut string) (*dflimg.Resource, error) {
	return a.db.FindResourceByShortcut(ctx, shortcut)
}
