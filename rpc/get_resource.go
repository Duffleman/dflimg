package rpc

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
)

// GetResource gets a resource and handles the response for it
func (r *RPC) GetResource(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	input := chi.URLParam(req, "input")

	resource, err := r.app.GetResource(ctx, input)
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
			err = tpl.Execute(w, map[string]interface{}{
				"resource": resource,
				"labels":   labelStr,
			})
			if err != nil {
				r.handleError(w, req, err)
			}
			return
		}
	}

	switch resource.Type {
	case "file":
		w.Header().Set("Content-Type", *resource.MimeType)

		b, modtime, err := r.app.GetS3File(ctx, resource)
		if err != nil {
			r.handleError(w, req, err)
			return
		}

		reader := bytes.NewReader(b)
		http.ServeContent(w, req, input, *modtime, reader)
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
