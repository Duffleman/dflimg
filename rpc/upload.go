package rpc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"dflimg/dflerr"
	"dflimg/rpc/middleware"
)

const maxUploadSize = 100 * 1024 // 100 MB
const uploadPath = "./tmp"

// Upload is an RPC handler for uploading a file
func (r *RPC) Upload(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	key := ctx.Value(middleware.UsernameKey)
	if key == nil || key == "" {
		r.handleError(w, req, dflerr.New(dflerr.AccessDenied, dflerr.M{"username": key}))
		return
	}

	labelsStr := req.FormValue("labels")
	var labels []string

	if labelsStr != "" {
		labels = strings.Split(labelsStr, ",")
	}

	file, _, err := req.FormFile("file")
	if err != nil {
		r.handleError(w, req, err)
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	io.Copy(&buf, file)

	res, err := r.app.Upload(ctx, buf, labels)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	accept := req.Header.Get("Accept")

	if strings.Contains(accept, "text/plain") {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(res.URL))
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}

	return
}
