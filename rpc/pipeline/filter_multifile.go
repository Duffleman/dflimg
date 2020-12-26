package pipeline

import (
	"dflimg/lib/cher"
)

// FilterMultiFile filters multi file requests so the rest of the chain doesn't
// need to worry.
func FilterMultiFile(p *Pipeline) (bool, error) {
	if p.context.multifile {
		return false, cher.New("multiple_files_unsupported", nil)
	}

	return true, nil
}
