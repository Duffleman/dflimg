package rpc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"dflimg"
	"dflimg/lib/cher"
	"dflimg/rpc/middleware"
)

// UploadFile is an RPC handler for uploading a file
func (r *RPC) UploadFile(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	key := ctx.Value(middleware.UsernameKey)
	if key == nil || key == "" {
		r.handleError(w, req, cher.New(cher.AccessDenied, cher.M{"username": key}))
		return
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		r.handleError(w, req, err)
		return
	}
	defer file.Close()

	var name = &header.Filename

	fileName := req.PostFormValue("name")
	if fileName != "" {
		name = &fileName
	}

	var buf bytes.Buffer
	io.Copy(&buf, file)

	res, err := r.app.UploadFile(ctx, &dflimg.CreateFileRequest{
		File: buf,
		Name: name,
	})
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
