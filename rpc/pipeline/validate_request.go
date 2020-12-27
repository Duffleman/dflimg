package pipeline

import (
	"dflimg/lib/cher"
)

// ValidateRequest validates some general request rules
func ValidateRequest(p *Pipeline) (bool, error) {
	if p.context.multifile {
		for _, i := range p.rwqs {
			rwq := i

			if rwq.r.Type == "url" {
				return false, cher.New("multiple_queries_with_url", nil)
			}
		}
	}

	return true, nil
}
