package app

import (
	"context"

	"dflimg"
)

func (a *App) ListLabels(ctx context.Context) ([]*dflimg.Label, error) {
	return a.db.ListLabels(ctx)
}
