package rpc

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"dflimg"
	"dflimg/app"
	"dflimg/lib/cher"
)

type handlerType func(*Pipeline) (bool, error)

// HandleResource loads a resource and handles the response for it
func (r *RPC) HandleResource(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	query := strings.TrimPrefix(req.URL.Path, "/")

	// remove the prefixed https://dfl.mn if present
	rootURL := dflimg.GetEnv("root_url") + "/"

	if strings.HasPrefix(query, rootURL) {
		query = strings.TrimPrefix(query, rootURL)
	}

	qi := r.app.ParseQueryType(query)

	resource, err := r.app.GetResource(ctx, qi)
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	pipe := &Pipeline{
		ctx:      ctx,
		app:      r.app,
		r:        req,
		w:        w,
		qi:       qi,
		resource: resource,
		steps: []handlerType{
			// ordered
			handleDeleted,
			handleNSFWPrimer,
			handleURLRedirect,
			loadFile,
			handleSyntaxHighlight,
			serveContent,
		},
	}

	err = pipe.Run()
	if err != nil {
		r.handleError(w, req, err)
		return
	}
}

type Pipeline struct {
	ctx      context.Context
	app      *app.App
	r        *http.Request
	w        http.ResponseWriter
	qi       *app.QueryInput
	resource *dflimg.Resource
	contents struct {
		modtime *time.Time
		bytes   []byte
	}
	steps []handlerType
}

func (p *Pipeline) Run() error {
	for _, fn := range p.steps {
		c, err := fn(p)
		if err != nil {
			return err
		}

		if !c {
			return nil
		}
	}

	return nil
}

func handleDeleted(p *Pipeline) (bool, error) {
	if p.resource.DeletedAt != nil {
		return false, cher.New(cher.NotFound, nil)
	}

	return true, nil
}

func handleNSFWPrimer(p *Pipeline) (bool, error) {
	// skip this step if the resource isn't NSFW
	if !p.resource.NSFW {
		return true, nil
	}

	// skip if we want to force a download
	if fd := p.r.URL.Query()["d"]; len(fd) >= 1 {
		return true, nil
	}

	// skip if we are already primed
	if _, ok := p.r.URL.Query()["primed"]; ok {
		return true, nil
	}

	tpl, err := template.ParseFiles("resources/nsfw.html")
	if err != nil {
		return false, err
	}
	err = tpl.Execute(p.w, map[string]interface{}{
		"resource": p.r,
	})

	return false, err
}

func handleURLRedirect(p *Pipeline) (bool, error) {
	// only handle URL types
	if p.resource.Type != "url" {
		return true, nil
	}

	p.w.Header().Set("Content-Type", "")

	http.Redirect(p.w, p.r, p.resource.Link, http.StatusTemporaryRedirect)

	return false, nil
}

func loadFile(p *Pipeline) (bool, error) {
	b, modtime, err := p.app.GetFile(p.ctx, p.resource)
	if err != nil {
		if c, ok := err.(cher.E); !ok || c.Code != cher.NotFound {
			return false, err
		}

		tpl, err := template.ParseFiles("resources/not_found.html")
		if err != nil {
			return false, err
		}

		err = tpl.Execute(p.w, nil)
		return false, err
	}

	p.contents.modtime = modtime
	p.contents.bytes = b

	return true, nil
}

func handleSyntaxHighlight(p *Pipeline) (bool, error) {
	accept := p.r.Header.Get("Accept")
	forceDownload := false

	if fd := p.r.URL.Query()["d"]; len(fd) >= 1 {
		forceDownload = true
	}

	isPlainText := strings.Contains(*p.resource.MimeType, "text/plain")
	acceptsHTML := strings.Contains(accept, "text/html")
	hasExt := p.qi.Ext != nil

	if forceDownload || !isPlainText || !acceptsHTML || !hasExt {
		return true, nil
	}

	tpl, err := template.ParseFiles("resources/code.html")
	if err != nil {
		return false, err
	}

	err = tpl.Execute(p.w, map[string]interface{}{
		"language": *p.qi.Ext,
		"title":    p.resource.Hash,
		"author":   p.resource.Owner,
		"content":  string(p.contents.bytes),
	})
	return false, err
}

func serveContent(p *Pipeline) (bool, error) {
	var display string = "inline"
	var mimetype string

	if p.resource.MimeType != nil {
		mimetype = *p.resource.MimeType
	}

	if fd := p.r.URL.Query()["d"]; len(fd) >= 1 {
		display = "attachment"
		mimetype = "application/octet-stream"
	}

	if p.resource.Name != nil {
		display = fmt.Sprintf("%s; filename=%s", display, *p.resource.Name)
	}

	p.w.Header().Set("Content-Type", mimetype)
	p.w.Header().Set("Content-Disposition", display)

	reader := bytes.NewReader(p.contents.bytes)

	http.ServeContent(p.w, p.r, p.qi.Filename(), *p.contents.modtime, reader)

	return true, nil
}
