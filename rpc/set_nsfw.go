package rpc

import (
	"encoding/json"
	"net/http"

	"dflimg"
	"dflimg/dflerr"
	"dflimg/rpc/middleware"
)

func (r *RPC) SetNSFW(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	key := ctx.Value(middleware.UsernameKey)
	if key == nil || key == "" {
		r.handleError(w, req, dflerr.New(dflerr.AccessDenied, dflerr.M{"username": key}))
		return
	}
	username := ctx.Value(middleware.UsernameKey).(string)

	body := &dflimg.SetNSFWRequest{}
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

	err = r.app.SetNSFW(ctx, resource.ID, body.NSFW)
	r.handleError(w, req, err)

	return
}
