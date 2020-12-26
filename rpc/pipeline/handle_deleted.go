package pipeline

import (
	"dflimg/lib/cher"
)

// HandleDeleted resources, show a 404 if not found
func HandleDeleted(p *Pipeline) (bool, error) {
	for _, i := range p.rwqs {
		rwq := i

		if rwq.r.DeletedAt != nil {
			return false, cher.New(cher.NotFound, nil)
		}
	}

	return true, nil
}
