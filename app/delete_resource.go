package app

import (
	"context"
	"errors"

	"dflimg"
	"dflimg/dflerr"
)

// DeleteResource deletes a resource
func (a *App) DeleteResource(ctx context.Context, resource *dflimg.Resource) error {
	if resource == nil {
		return errors.New("empty resource given")
	}

	if resource.DeletedAt != nil {
		return dflerr.New("already_deleted", nil)
	}

	return a.db.DeleteResource(ctx, resource.ID)
}
