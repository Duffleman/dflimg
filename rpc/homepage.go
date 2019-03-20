package rpc

import (
	"net/http"
)

// Homepage is the root page for the site
func (r *RPC) Homepage(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "") // Needed for redirect to work
	http.Redirect(w, req, "https://duffleman.co.uk", http.StatusMovedPermanently)

	return
}
