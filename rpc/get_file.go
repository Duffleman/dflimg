package rpc

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (r *RPC) GetFile(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	fileID := chi.URLParam(req, "fileID")

	fileType, res, err := r.app.GetFile(ctx, fileID)
	if err != nil {
		r.handleError(w, req, err, &map[string]interface{}{
			"fileID": fileID,
		})
		return
	}

	w.Header().Set("Content-Type", fileType)
	w.Write(res.Bytes())

	return
}
