package rpc

import (
	"net/http"
	"strings"

	"dflimg"
	"dflimg/dflerr"
	"dflimg/rpc/middleware"
)

func (r *RPC) TagResource(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	user := ctx.Value(middleware.UsernameKey)
	if user == nil || user == "" {
		r.handleError(w, req, dflerr.New(dflerr.AccessDenied, dflerr.M{"username": user}))
		return
	}

	urlStr := req.FormValue("url")
	if urlStr == "" {
		r.handleError(w, req, dflerr.New(dflerr.RequestFailure, dflerr.M{"missing_key": "url"}))
		return
	}

	tagsStr := req.FormValue("tags")
	if tagsStr == "" {
		r.handleError(w, req, dflerr.New(dflerr.RequestFailure, dflerr.M{"missing_key": "tags"}))
		return
	}

	tags := strings.Split(tagsStr, ",")

	entry := urlStr
	rootURL := dflimg.GetEnv("root_url") + "/"

	if strings.HasPrefix(urlStr, rootURL) {
		entry = strings.TrimPrefix(urlStr, rootURL)
	}

	resource, err := r.app.GetResource(ctx, entry)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	if resource.Owner != user {
		r.handleError(w, req, dflerr.New(dflerr.AccessDenied, nil))
		return
	}

	err = r.app.TagResource(ctx, resource.ID, tags)
	if err != nil {
		r.handleError(w, req, err)
	}

	return
}
