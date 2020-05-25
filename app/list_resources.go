package app

import (
	"context"

	"dflimg"
)

// ListResources returns a list of all resources for a user
func (a *App) ListResources(ctx context.Context, req *dflimg.ListResourcesRequest) ([]*dflimg.Resource, error) {
	return a.db.ListResources(ctx, req.Username, req.IncludeDeleted)
}
