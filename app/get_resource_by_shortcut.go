package app

import (
	"context"

	"dflimg"

	"golang.org/x/sync/errgroup"
)

// GetResourceByShortcut returns a resource from a :shortcut
func (a *App) GetResourceByShortcut(ctx context.Context, shortcut string) (*dflimg.Resource, error) {
	var resource *dflimg.Resource
	var labels []*dflimg.Label

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() (err error) {
		resource, err = a.db.FindResourceByShortcut(gctx, shortcut)
		return err
	})

	g.Go(func() (err error) {
		labels, err = a.db.GetLabelsByShortcut(gctx, shortcut)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	resource.Labels = labels

	return resource, nil
}
