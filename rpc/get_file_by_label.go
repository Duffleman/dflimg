package rpc

import (
	"dflimg"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

func (r *RPC) GetFileByLabel(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	label := chi.URLParam(req, "label")

	res, err := r.app.GetFileByLabel(ctx, label)
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
