package rpc

import (
	"encoding/json"
	"net/http"

	"dflimg"
	"dflimg/lib/cher"
	"dflimg/rpc/middleware"
)

func (r *RPC) ListResources(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	key := ctx.Value(middleware.UsernameKey)
	if key == nil || key == "" {
		r.handleError(w, req, cher.New(cher.AccessDenied, cher.M{"username": key}))
		return
	}
	username := ctx.Value(middleware.UsernameKey).(string)

	body := &dflimg.ListResourcesRequest{}
	err := json.NewDecoder(req.Body).Decode(body)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	if body.Username == nil || *body.Username != username {
		r.handleError(w, req, cher.New(cher.AccessDenied, nil))
		return
	}

	resources, err := r.app.ListResources(ctx, body)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	json.NewEncoder(w).Encode(resources)

	return
}
