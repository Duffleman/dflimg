package rpc

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"dflimg/dflerr"
	"dflimg/rpc/middleware"
)

func (r *RPC) ShortenURL(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	key := ctx.Value(middleware.UsernameKey)
	if key == nil || key == "" {
		r.handleError(w, req, dflerr.New(dflerr.AccessDenied, dflerr.M{"username": key}))
		return
	}

	shortcutsStr := req.FormValue("shortcuts")
	urlStr := req.FormValue("url")
	var shortcuts []string

	if urlStr == "" {
		err := errors.New("missing url")
		r.handleError(w, req, err)
		return
	}

	if shortcutsStr != "" {
		shortcuts = strings.Split(shortcutsStr, ",")
	}

	res, err := r.app.ShortenURL(ctx, urlStr, shortcuts)
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
