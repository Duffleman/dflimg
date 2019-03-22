package rpc

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (r *RPC) GetFileByLabel(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	label := chi.URLParam(req, "label")

	fileType, res, err := r.app.GetFileByLabel(ctx, label)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	w.Header().Set("Content-Type", fileType)
	w.Write(res.Bytes())

	return
}
