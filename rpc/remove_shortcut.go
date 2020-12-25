package rpc

import (
	"encoding/json"
	"net/http"

	"dflimg"
	"dflimg/app"
	"dflimg/lib/cher"
	"dflimg/rpc/middleware"
)

func (r *RPC) RemoveShortcut(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	key := ctx.Value(middleware.UsernameKey)
	if key == nil || key == "" {
		r.handleError(w, req, cher.New(cher.AccessDenied, nil))
		return
	}
	username := ctx.Value(middleware.UsernameKey).(string)

	body := &dflimg.ChangeShortcutRequest{}
	err := json.NewDecoder(req.Body).Decode(body)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	qi := r.app.ParseQueryType(body.Query)

	if qi.QueryType == app.Name {
		r.handleError(w, req, cher.New("cannot_query_resource_by_name", cher.M{"query": qi}))
		return
	}

	resource, err := r.app.GetResource(ctx, qi)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	if resource.Owner != username {
		r.handleError(w, req, cher.New(cher.AccessDenied, nil))
		return
	}

	err = r.app.RemoveShortcut(ctx, resource, body.Shortcut)
	r.handleError(w, req, err)

	return
}
