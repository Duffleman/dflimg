package rpc

import (
	"encoding/json"
	"net/http"

	"dflimg"
	"dflimg/app"
	"dflimg/lib/cher"
	"dflimg/rpc/middleware"
)

func (r *RPC) DeleteResource(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	key := ctx.Value(middleware.UsernameKey)
	if key == nil || key == "" {
		r.handleError(w, req, cher.New(cher.AccessDenied, nil))
		return
	}
	username := ctx.Value(middleware.UsernameKey).(string)

	body := &dflimg.IdentifyResource{}
	err := json.NewDecoder(req.Body).Decode(body)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	qi := r.app.ParseQueryType(body.Query)

	if len(qi) != 1 {
		r.handleError(w, req, cher.New("multi_query_not_supported", cher.M{"query": qi}))
		return
	}

	if qi[0].QueryType == app.Name {
		r.handleError(w, req, cher.New("cannot_query_resource_by_name", cher.M{"query": qi}))
		return
	}

	resource, err := r.app.GetResource(ctx, qi[0])
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	if resource.Owner != username {
		r.handleError(w, req, cher.New(cher.AccessDenied, nil))
		return
	}

	err = r.app.DeleteResource(ctx, resource)
	r.handleError(w, req, err)

	return
}
