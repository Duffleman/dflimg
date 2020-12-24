package rpc

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"dflimg"
	"dflimg/dflerr"

	"github.com/go-chi/chi"
)

// GetResource gets a resource and handles the response for it
func (r *RPC) GetResource(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var forceDownload bool

	query := chi.URLParam(req, "query")
	accept := req.Header.Get("Accept")

	if fd := req.URL.Query()["d"]; len(fd) >= 1 {
		forceDownload = true
	}

	resource, ext, err := r.app.GetResource(ctx, query)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	if resource.DeletedAt != nil {
		r.handleError(w, req, dflerr.New(dflerr.NotFound, nil))
		return
	}

	if resource.NSFW {
		if _, ok := req.URL.Query()["primed"]; !ok {
			tpl, err := template.ParseFiles("resources/nsfw.html")
			if err != nil {
				r.handleError(w, req, err)
				return
			}
			err = tpl.Execute(w, map[string]interface{}{
				"resource": resource,
			})
			if err != nil {
				r.handleError(w, req, err)
			}
			return
		}
	}

	switch resource.Type {
	case "file":
		b, modtime, err := r.app.GetFile(ctx, resource)
		if err != nil {
			if err == dflerr.ErrNotFound {
				tpl, err := template.ParseFiles("resources/not_found.html")
				if err != nil {
					r.handleError(w, req, err)
					return
				}
				err = tpl.Execute(w, nil)
				if err != nil {
					r.handleError(w, req, err)
				}
				return
			}

			r.handleError(w, req, err)
			return
		}

		reader := bytes.NewReader(b)

		isPlainText := strings.Contains(*resource.MimeType, "text/plain")
		acceptsHTML := strings.Contains(accept, "text/html")
		hasNoExt := ext == nil

		if isPlainText && acceptsHTML && !hasNoExt {
			// do the formatting
			tpl, err := template.ParseFiles("resources/code.html")
			if err != nil {
				r.handleError(w, req, err)
				return
			}

			language := *resource.MimeType

			if ext != nil {
				language = *ext
			}

			err = tpl.Execute(w, map[string]interface{}{
				"language": language,
				"title":    resource.Hash,
				"author":   resource.Owner,
				"content":  string(b),
			})
			if err != nil {
				r.handleError(w, req, err)
			}
			return
		}

		writeHeaders(w, resource, forceDownload)
		http.ServeContent(w, req, query, *modtime, reader)
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

func writeHeaders(w http.ResponseWriter, resource *dflimg.Resource, forceDownload bool) {
	var display string = "inline"
	var mimetype string

	if resource.Type == "url" {
		w.Header().Set("Content-Type", "")
		return
	}

	if resource.MimeType != nil {
		mimetype = *resource.MimeType
	}

	if forceDownload {
		display = "attachment"
		mimetype = "application/octet-stream"
	}

	if resource.Name != nil {
		display = fmt.Sprintf("%s; filename=%s", display, *resource.Name)
	}

	w.Header().Set("Content-Type", mimetype)
	w.Header().Set("Content-Disposition", display)
}
