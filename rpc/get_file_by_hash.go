package rpc

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (r *RPC) GetFileByHash(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	hash := chi.URLParam(req, "hash")

	fileType, res, err := r.app.GetFileByHash(ctx, hash)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	w.Header().Set("Content-Type", fileType)
	w.Write(res.Bytes())

	return
}
