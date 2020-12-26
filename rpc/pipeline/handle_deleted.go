package pipeline

import (
	"dflimg/lib/cher"
)

func HandleDeleted(p *Pipeline) (bool, error) {
	for _, rwq := range p.rwqs {
		if rwq.r.DeletedAt != nil {
			return false, cher.New(cher.NotFound, nil)
		}
	}

	return true, nil
}
