package app

import (
	"context"

	"dflimg"
	"dflimg/dflerr"

	"golang.org/x/sync/errgroup"
)

func (a *App) GetResourceByHash(ctx context.Context, hash string) (*dflimg.Resource, error) {
	serial, err := a.decodeHash(hash)
	if err != nil {
		return nil, err
	}

	var resource *dflimg.Resource
	var labels []*dflimg.Label

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() (err error) {
		resource, err = a.db.FindResourceBySerial(gctx, serial)
		return err
	})

	g.Go(func() (err error) {
		labels, err = a.db.GetLabelsBySerial(gctx, serial)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	resource.Labels = labels

	return resource, nil
}

func (a *App) decodeHash(hash string) (int, error) {
	var set []int

	set, err := a.hasher.DecodeWithError(hash)
	if len(set) != 1 {
		return 0, dflerr.New("cannot decode hash", dflerr.M{"hash": hash}, dflerr.New("expecting single hashed item in body", nil))
	}

	return set[0], err
}
