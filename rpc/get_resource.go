package rpc

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"dflimg"
	"dflimg/dflerr"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/go-chi/chi"
)

// GetResource gets a resource and handles the response for it
func (r *RPC) GetResource(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	query := chi.URLParam(req, "query")
	accept := req.Header.Get("Accept")

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
		b, modtime, err := r.app.GetS3File(ctx, resource)
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

		// Let's try to format the output if possible
		// - you must accept text/html
		// - it must be a text document
		if strings.Contains(*resource.MimeType, "text/plain") && strings.Contains(accept, "text/html") {
			var lexer chroma.Lexer

			// if you provide an extension, format the doc accordingly
			if ext != nil {
				// if you're looking for .txt, keep it to text, no HTML
				if *ext == "txt" {
					fallback(resource, w, req, query, *modtime, reader)
					return
				} else {
					// match the lexer to the extension given
					lexer = lexers.Match(fmt.Sprintf("file.%s", *ext))
				}
			} else {
				// analyse the document to figure out what it may be
				lexer = lexers.Analyse(string(b))
			}

			if lexer == nil {
				lexer = lexers.Fallback
			}

			lexer = chroma.Coalesce(lexer)

			formatter := formatters.Get("html")
			if formatter == nil {
				formatter = formatters.Fallback
			}

			style := styles.Get("vs")
			if style == nil {
				style = styles.Fallback
			}

			contents, err := ioutil.ReadAll(reader)
			if err != nil {
				fallback(resource, w, req, query, *modtime, reader)
				return
			}

			iterator, err := lexer.Tokenise(nil, string(contents))
			if err != nil {
				fallback(resource, w, req, query, *modtime, reader)
				return
			}

			err = formatter.Format(w, style, iterator)
			if err != nil {
				fallback(resource, w, req, query, *modtime, reader)
				return
			}

			return
		}

		fallback(resource, w, req, query, *modtime, reader)
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

func fallback(resource *dflimg.Resource, w http.ResponseWriter, req *http.Request, query string, modtime time.Time, reader io.ReadSeeker) {
	w.Header().Set("Content-Type", *resource.MimeType)

	http.ServeContent(w, req, query, modtime, reader)

	return
}
