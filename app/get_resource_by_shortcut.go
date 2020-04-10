package app

import (
	"context"

	"dflimg"
)

// GetResourceByShortcut returns a resource from a :shortcut
func (a *App) GetResourceByShortcut(ctx context.Context, shortcut string) (*dflimg.Resource, error) {
	return a.db.FindResourceByShortcut(ctx, shortcut)
}
