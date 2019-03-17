package rpc

import (
	"dflimg"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

func (r *RPC) GetFile(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	fileID := chi.URLParam(req, "fileID")

	res, err := r.app.GetFile(ctx, fileID)
	if err != nil {
		r.logger.WithError(err)
		if err == dflimg.ErrNotFound {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(map[string]string{"code": "not_found"})
		} else {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"code": err.Error()})
		}
		return
	}

	w.Write(res.Bytes())

	return
}
