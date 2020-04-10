package rpc

import (
	"encoding/json"
	"net/http"

	"dflimg"
	"dflimg/dflerr"
	"dflimg/rpc/middleware"
)

func (r *RPC) ViewDetails(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

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

	if resource.DeletedAt != nil {
		key := ctx.Value(middleware.UsernameKey)
		if key == nil || key == "" {
			r.handleError(w, req, dflerr.New(dflerr.AccessDenied, nil))
			return
		}
		username := ctx.Value(middleware.UsernameKey).(string)

		if resource.Owner != username {
			r.handleError(w, req, dflerr.New(dflerr.NotFound, nil))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)

	return
}
