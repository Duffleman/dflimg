package app

import (
	"context"
	"dflimg/dflerr"
)

func (a *App) TagResource(ctx context.Context, resourceID string, tags []string) error {
	labels, err := a.db.GetLabelsByName(ctx, tags)
	if err != nil {
		return err
	}

	if len(labels) != len(tags) {
		return dflerr.New(dflerr.NotFound, nil)
	}

	return a.db.TagResource(ctx, resourceID, labels)
}
