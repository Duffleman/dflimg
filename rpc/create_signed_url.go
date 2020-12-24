package rpc

import (
	"encoding/json"
	"net/http"
	"strings"

	"dflimg"
	"dflimg/lib/cher"
	"dflimg/rpc/middleware"
)

// CreateSignedURL creates a signed URL for file uploads
func (r *RPC) CreateSignedURL(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	key := ctx.Value(middleware.UsernameKey)
	if key == nil || key == "" {
		r.handleError(w, req, cher.New(cher.AccessDenied, cher.M{"username": key}))
		return
	}
	username := ctx.Value(middleware.UsernameKey).(string)

	body := &dflimg.CreateSignedURLRequest{}
	err := json.NewDecoder(req.Body).Decode(body)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	res, err := r.app.CreateSignedURL(ctx, username, body.Name, body.ContentType)
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
