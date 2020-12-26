package pipeline

import (
	"net/http"
)

// HandleURLRedirect handles URL redirects
func HandleURLRedirect(p *Pipeline) (bool, error) {
	if p.context.multifile || p.rwqs[0].r.Type != "url" {
		return true, nil
	}

	p.w.Header().Set("Content-Type", "")

	http.Redirect(p.w, p.r, p.rwqs[0].r.Link, http.StatusTemporaryRedirect)

	return false, nil
}
