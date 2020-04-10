package app

import (
	"context"

	"dflimg"
)

// GetResourceByHash returns a resource when given a hash
func (a *App) GetResourceByHash(ctx context.Context, hash string) (*dflimg.Resource, error) {
	return a.db.FindResourceByHash(ctx, hash)
}
