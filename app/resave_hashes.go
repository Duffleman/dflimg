package app

import (
	"context"
	"sync"

	"dflimg"
)

// ResaveHashes finds all resources without a saved hash and saves it to the DB
func (a *App) ResaveHashes(ctx context.Context) error {
	resources, err := a.db.ListResourcesWithoutHash(ctx)
	if err != nil {
		return err
	}
	errCh := make(chan error, 1)
	c := make(chan *dflimg.ShortFormResource)
	wg := &sync.WaitGroup{}

	// I moved this above to stop blocking execution.
	// as we're now looking for an empty channel (Exhibit A) it would have been
	// empty immediately had it happened after the workers
	go func() {
		for _, r := range resources {
			c <- r
		}
		close(c)
	}()

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go a.doDaWork(ctx, wg, c, errCh)
	}

	if len(errCh) > 0 {
		err := <-errCh
		return err
	}

	wg.Wait()

	return nil
}

func (a *App) doDaWork(ctx context.Context, wg *sync.WaitGroup, ch chan *dflimg.ShortFormResource, errCh chan error) {
	defer wg.Done()

	for {
		select {
		case r, ok := <-ch:
			if !ok {
				// Exhibit A: the channel is now closed and empty
				return
			}

			hash := a.makeHash(r.Serial)
			// if SaveHash returns an error, whack it on the channel
			err := a.db.SaveHash(ctx, r.Serial, hash)
			if err != nil {
				errCh <- err
				return
			}

		case <-ctx.Done():
			// ctx is cancelled, throw the error on the channel
			errCh <- ctx.Err()
			return
		}
	}
}
