package app

import (
	"context"

	"dflimg"
	"dflimg/dflerr"
)

// DeleteResource deletes a resource
func (a *App) DeleteResource(ctx context.Context, resource *dflimg.Resource) error {
	if resource == nil {
		return dflerr.New(dflerr.NotFound, nil)
	}

	if resource.DeletedAt != nil {
		return dflerr.New(dflerr.NotFound, nil)
	}

	return a.db.DeleteResource(ctx, resource.ID)
}
