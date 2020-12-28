package pipeline

import (
	"dflimg"
	"sync"

	"golang.org/x/sync/errgroup"
)

// ListFromDB loads all queries and finds matching resources in the DB. It will
// maintain the exact order too!
func ListFromDB(p *Pipeline) (bool, error) {
	g, gctx := errgroup.WithContext(p.ctx)

	rm := &resourceMap{
		s: make(map[string]*dflimg.Resource),
	}

	for _, i := range p.qi {
		query := i

		g.Go(func() (err error) {
			resource, err := p.app.GetResource(gctx, query)
			if err != nil {
				return err
			}

			rm.Add(query.Original, resource)

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return false, err
	}

	var rwqs []*resourceWithQuery

	for _, i := range p.qi {
		query := i

		rwqs = append(rwqs, &resourceWithQuery{
			qi: query,
			r:  rm.s[query.Original],
		})
	}

	p.rwqs = rwqs

	return true, nil
}

type resourceMap struct {
	s map[string]*dflimg.Resource
	sync.Mutex
}

func (rm *resourceMap) Add(key string, r *dflimg.Resource) {
	rm.Lock()
	defer rm.Unlock()

	rm.s[key] = r
}
