package app

import (
	"context"

	"dflimg"

	"golang.org/x/sync/errgroup"
)

// GetResourceByHash returns a resource when given a hash
func (a *App) GetResourceByHash(ctx context.Context, hash string) (*dflimg.Resource, error) {
	var resource *dflimg.Resource
	var labels []*dflimg.Label

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() (err error) {
		resource, err = a.db.FindResourceByHash(gctx, hash)
		return err
	})

	g.Go(func() (err error) {
		labels, err = a.db.GetLabelsByHash(gctx, hash)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	resource.Labels = labels

	return resource, nil
}
