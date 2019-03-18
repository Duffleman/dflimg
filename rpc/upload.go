package rpc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

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

	labelsStr := req.FormValue("labels")
	var labels []string

	if labelsStr != "" {
		labels = strings.Split(labelsStr, ",")
	}

	file, _, err := req.FormFile("file")
	if err != nil {
		r.logger.WithError(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"code": err.Error()})
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	io.Copy(&buf, file)

	res, err := r.app.Upload(ctx, buf, labels)
	if err != nil {
		r.logger.WithError(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"code": err.Error()})
		return
	}

	accept := req.Header.Get("Accept")

	if accept == "text/plain" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(res.URL))
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}

	return
}
