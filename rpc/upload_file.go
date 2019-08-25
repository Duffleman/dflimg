package rpc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"dflimg"
	"dflimg/dflerr"
	"dflimg/rpc/middleware"
)

// UploadFile is an RPC handler for uploading a file
func (r *RPC) UploadFile(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	key := ctx.Value(middleware.UsernameKey)
	if key == nil || key == "" {
		r.handleError(w, req, dflerr.New(dflerr.AccessDenied, dflerr.M{"username": key}))
		return
	}

	shortcutsStr := req.FormValue("shortcuts")
	nsfwStr := req.FormValue("nsfw")
	var shortcuts []string
	var nsfw bool

	switch nsfwStr {
	case "true":
		nsfw = true
	default:
		nsfw = false
	}

	if shortcutsStr != "" {
		shortcuts = strings.Split(shortcutsStr, ",")
	}

	file, _, err := req.FormFile("file")
	if err != nil {
		r.handleError(w, req, err)
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	io.Copy(&buf, file)

	resourceReq := &dflimg.CreateResourceRequest{
		Type:      "file",
		File:      buf,
		Shortcuts: shortcuts,
		NSFW:      nsfw,
	}

	res, err := r.app.UploadFile(ctx, resourceReq)
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
