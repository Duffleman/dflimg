package rpc

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (r *RPC) GetFileByShortcut(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	shortcut := chi.URLParam(req, "shortcut")

	fileType, res, err := r.app.GetFileByShortcut(ctx, shortcut)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	w.Header().Set("Content-Type", fileType)
	w.Write(res.Bytes())

	return
}
