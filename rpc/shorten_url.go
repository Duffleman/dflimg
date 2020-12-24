package rpc

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"dflimg"
	"dflimg/lib/cher"
	"dflimg/rpc/middleware"
)

func (r *RPC) ShortenURL(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	key := ctx.Value(middleware.UsernameKey)
	if key == nil || key == "" {
		r.handleError(w, req, cher.New(cher.AccessDenied, cher.M{"username": key}))
		return
	}

	body := &dflimg.CreateURLRequest{}
	err := json.NewDecoder(req.Body).Decode(body)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	if body.URL == "" {
		err := errors.New("missing_url")
		r.handleError(w, req, err)
		return
	}

	res, err := r.app.ShortenURL(ctx, body.URL)
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
