package rpc

import (
	"dflimg"
	"encoding/json"
	"net/http"

	"dflimg/dflerr"
	"dflimg/rpc/middleware"
)

func (r *RPC) DeleteResource(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	key := ctx.Value(middleware.UsernameKey)
	if key == nil || key == "" {
		r.handleError(w, req, dflerr.New(dflerr.AccessDenied, nil))
		return
	}
	username := ctx.Value(middleware.UsernameKey).(string)

	body := &dflimg.IdentifyResource{}
	err := json.NewDecoder(req.Body).Decode(body)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	resource, err := r.app.GetResource(ctx, body.Query)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	if resource.Owner != username {
		r.handleError(w, req, dflerr.New(dflerr.AccessDenied, nil))
		return
	}

	err = r.app.DeleteResource(ctx, resource)
	r.handleError(w, req, err)

	return
}
