package app

import (
	"context"

	"dflimg"
	"dflimg/lib/cher"
)

// DeleteResource deletes a resource
func (a *App) DeleteResource(ctx context.Context, resource *dflimg.Resource) error {
	if resource == nil {
		return cher.New(cher.NotFound, nil)
	}

	if resource.DeletedAt != nil {
		return cher.New(cher.NotFound, nil)
	}

	return a.db.DeleteResource(ctx, resource.ID)
}
