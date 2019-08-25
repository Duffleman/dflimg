package rpc

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"dflimg"

	"github.com/go-chi/chi"
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

	if resource.NSFW {
		if _, ok := req.URL.Query()["primed"]; !ok {
			labelStr := resource.StringifyLabels()

			tpl, err := template.ParseFiles("resources/nsfw.html")
			if err != nil {
				r.handleError(w, req, err)
				return
			}
			tpl.Execute(w, map[string]interface{}{
				"resource": resource,
				"labels":   labelStr,
			})
			return
		}
	}

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