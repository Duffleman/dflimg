package rpc

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"dflimg/dflerr"
	"dflimg/rpc/middleware"
)

// CreatedSignedURL creates a signed URL for file uploads
func (r *RPC) CreatedSignedURL(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	key := ctx.Value(middleware.UsernameKey)
	if key == nil || key == "" {
		r.handleError(w, req, dflerr.New(dflerr.AccessDenied, dflerr.M{"username": key}))
		return
	}
	username := ctx.Value(middleware.UsernameKey).(string)

	contentType := req.FormValue("content-type")
	contentLengthStr := req.FormValue("content-length")
	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	res, err := r.app.CreatedSignedURL(ctx, username, contentType, contentLength, nil, false)
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
