package pipeline

import (
	"net/http"
)

func HandleURLRedirect(p *Pipeline) (bool, error) {
	// only handle URL types
	if p.context.multifile || p.rwqs[0].r.Type != "url" {
		return true, nil
	}

	p.w.Header().Set("Content-Type", "")

	http.Redirect(p.w, p.r, p.rwqs[0].r.Link, http.StatusTemporaryRedirect)

	return false, nil
}
