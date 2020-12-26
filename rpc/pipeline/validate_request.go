package pipeline

import (
	"dflimg/lib/cher"
)

func ValidateRequest(p *Pipeline) (bool, error) {
	if p.context.multifile {
		for _, rwq := range p.rwqs {
			if rwq.r.Type == "url" {
				return false, cher.New("multi_file_with_url", nil)
			}
		}
	}

	return true, nil
}
