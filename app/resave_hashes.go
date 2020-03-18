package app

import (
	"context"

	"dflimg"
)

func (a *App) ResaveHashes(ctx context.Context) error {
	resources, err := a.db.ListResourcesWithoutHash(ctx)
	if err != nil {
		return err
	}

	c := make(chan *dflimg.ShortFormResource)

	// start 3 workers
	go a.doDaWork(ctx, c)
	go a.doDaWork(ctx, c)
	go a.doDaWork(ctx, c)

	// dump all resources into the channel
	for _, r := range resources {
		c <- r
	}

	return nil
}

func (a *App) doDaWork(ctx context.Context, ch chan *dflimg.ShortFormResource) {
	for {
		select {
		case r := <-ch:
			hash := a.makeHash(r.Serial)

			a.db.SaveHash(ctx, r.Serial, hash)
		}
	}
}
