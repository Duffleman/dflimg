package app

import (
	"context"
	"dflimg/dflerr"
)

// TagResource tags a resource after it has been uploaded
func (a *App) TagResource(ctx context.Context, resourceID string, tags []string, nsfw bool) error {
	labels, err := a.db.GetLabelsByName(ctx, tags)
	if err != nil {
		return err
	}

	if len(labels) != len(tags) {
		return dflerr.New(dflerr.NotFound, nil)
	}

	err = a.db.SetNSFW(ctx, resourceID, nsfw)
	if err != nil {
		return err
	}

	return a.db.TagResource(ctx, resourceID, labels)
}
