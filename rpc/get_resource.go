package rpc

import (
	"errors"
	"net/http"
	"strings"

	"dflimg"

	"github.com/go-chi/chi"
	"github.com/kr/pretty"
)

// ShortcutCharacter marks the character used to find shortcuts
const ShortcutCharacter = ":"

// GetResource gets a resource and handles the response for it
func (r *RPC) GetResource(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	input := chi.URLParam(req, "input")

	var resource *dflimg.Resource
	var err error

	if strings.HasPrefix(input, ShortcutCharacter) {
		resource, err = r.app.GetResourceByShortcut(ctx, input)
	} else {
		resource, err = r.app.GetResourceByHash(ctx, input)
	}
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	pretty.Println(resource)

	switch resource.Type {
	case "file":
		w.Header().Set("Content-Type", *resource.MimeType)

		bytes, err := r.app.GetS3File(ctx, resource)
		if err != nil {
			r.handleError(w, req, err)
			return
		}

		w.Write(bytes)
		return
	case "url":
		w.Header().Set("Content-Type", "") // Needed for redirect to work
		http.Redirect(w, req, resource.Link, http.StatusTemporaryRedirect)
		return
	default:
		r.handleError(w, req, errors.New("unknown resource type"))
		return
	}

}
