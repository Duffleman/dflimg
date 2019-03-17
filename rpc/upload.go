package rpc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"dflimg/rpc/middleware"
)

const maxUploadSize = 100 * 1024 // 100 MB
const uploadPath = "./tmp"

// Upload is an RPC handler for uploading a file
func (r *RPC) Upload(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	key := ctx.Value(middleware.UsernameKey)
	if key == nil || key == "" {
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(map[string]string{"code": "access_denied"})
		return
	}

	file, _, err := req.FormFile("file")
	if err != nil {
		r.logger.WithError(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	io.Copy(&buf, file)

	res, err := r.app.Upload(ctx, buf)
	if err != nil {
		r.logger.WithError(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(res)

	return
}
