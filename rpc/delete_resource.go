package rpc

import (
	"dflimg"
	"net/http"
	"strings"

	"dflimg/dflerr"
	"dflimg/rpc/middleware"
)

func (r *RPC) DeleteResource(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	username := ctx.Value(middleware.UsernameKey)
	if username == nil || username == "" {
		r.handleError(w, req, dflerr.New(dflerr.AccessDenied, nil))
		return
	}

	input := req.FormValue("input")
	rootURL := dflimg.GetEnv("root_url") + "/"

	if strings.HasPrefix(input, rootURL) {
		input = strings.TrimPrefix(input, rootURL)
	}

	resource, err := r.app.GetResource(ctx, input)
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
