package rpc

import (
	"encoding/json"
	"net/http"
)

func (r *RPC) ListLabels(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	labels, err := r.app.ListLabels(ctx)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(labels)
}
