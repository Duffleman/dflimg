package rpc

import (
	"net/http"

	"dflimg/dflerr"
	"dflimg/rpc/middleware"
)

func (r *RPC) ResaveHashes(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	key := ctx.Value(middleware.UsernameKey)
	if key == nil || key == "" {
		r.handleError(w, req, dflerr.New(dflerr.AccessDenied, nil))
		return
	}

	err := r.app.ResaveHashes(ctx)
	r.handleError(w, req, err)

	return
}
