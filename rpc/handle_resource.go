package rpc

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	tt "text/template"
	"time"

	"dflimg"
	"dflimg/app"
	"dflimg/lib/cher"

	"github.com/gomarkdown/markdown"
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
		context:  make(map[string]bool),
		steps: []handlerType{
			// ordered
			handleDeleted,
			makeContext,
			handleNSFWPrimer,
			handleURLRedirect,
			loadFile,
			handleMdToHTML,
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
	context map[string]bool
	steps   []handlerType
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

func makeContext(p *Pipeline) (bool, error) {
	if fd := p.r.URL.Query()["d"]; len(fd) >= 1 {
		p.context["forceDownload"] = true
	}

	if _, ok := p.r.URL.Query()["primed"]; ok {
		p.context["primed"] = true
	}

	if p.resource.MimeType != nil && strings.Contains(*p.resource.MimeType, "text/plain") {
		p.context["resourceIsText"] = true
	}

	if strings.Contains(p.r.Header.Get("Accept"), "text/html") {
		p.context["acceptsHTML"] = true
	}

	return true, nil
}

func handleNSFWPrimer(p *Pipeline) (bool, error) {
	// skip this step if qualifiers are met
	if !p.resource.NSFW || p.context["forceDownload"] || p.context["primed"] {
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

func handleMdToHTML(p *Pipeline) (bool, error) {
	resourceIsText := p.context["resourceIsText"]
	acceptsHTML := p.context["acceptsHTML"]
	hasExt := p.qi.Ext != nil && *p.qi.Ext == "mdhtml"

	// skip if we don't meet the qualifier
	if p.qi.QueryType == app.Name || !hasExt || !acceptsHTML || !resourceIsText {
		return true, nil
	}

	output := markdown.ToHTML(p.contents.bytes, nil, nil)

	display, _ := p.getContentHeaders()
	mimetype := "text/html; charset=utf-8"

	if fd := p.r.URL.Query()["d"]; len(fd) >= 1 {
		mimetype = "application/octet-stream"
	}

	p.w.Header().Set("Content-Type", mimetype)
	p.w.Header().Set("Content-Disposition", display)

	tpl, err := tt.ParseFiles("resources/markdown.html")
	if err != nil {
		return false, err
	}

	title := p.resource.Hash

	if p.resource.Name != nil {
		title = p.resource.Name
	}

	err = tpl.Execute(p.w, map[string]interface{}{
		"title":   title,
		"author":  p.resource.Owner,
		"content": string(output),
	})

	return false, err
}

func handleSyntaxHighlight(p *Pipeline) (bool, error) {
	resourceIsText := p.context["resourceIsText"]
	acceptsHTML := p.context["acceptsHTML"]
	forceDownload := p.context["forceDownload"]
	hasExt := p.qi.Ext != nil

	if p.qi.QueryType == app.Name || forceDownload || !resourceIsText || !acceptsHTML || !hasExt {
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
	display, mimetype := p.getContentHeaders()

	p.w.Header().Set("Content-Type", mimetype)
	p.w.Header().Set("Content-Disposition", display)

	reader := bytes.NewReader(p.contents.bytes)

	http.ServeContent(p.w, p.r, p.qi.Filename(), *p.contents.modtime, reader)

	return true, nil
}

func (p *Pipeline) getContentHeaders() (string, string) {
	var display string = "inline"
	var mimetype string

	if p.resource.MimeType != nil {
		mimetype = *p.resource.MimeType
	}

	if p.context["forceDownload"] {
		display = "attachment"
		mimetype = "application/octet-stream"
	}

	if p.resource.Name != nil {
		display = fmt.Sprintf("%s; filename=%s", display, *p.resource.Name)
	}

	return display, mimetype
}
