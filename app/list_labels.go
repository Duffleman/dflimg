package app

import (
	"context"

	"dflimg"
)

// ListLabels returns the list of labels you can use
func (a *App) ListLabels(ctx context.Context) ([]*dflimg.Label, error) {
	return a.db.ListLabels(ctx)
}
