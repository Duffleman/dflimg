package app

import (
	"context"

	"dflimg"
)

// GetResourceByName returns a resource from it's name
func (a *App) GetResourceByName(ctx context.Context, name string) (*dflimg.Resource, error) {
	return a.db.FindResourceByName(ctx, name)
}
