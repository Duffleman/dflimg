package pipeline

import (
	"fmt"
	"strings"
	tt "text/template"

	"github.com/gomarkdown/markdown"
)

// RenderMD renders a file set from Markdown to HTMLs
func RenderMD(p *Pipeline) (bool, error) {
	allAreText := true

	for _, i := range p.rwqs {
		rwq := i

		if !rwq.context.isText {
			allAreText = false
		}

	}

	switch {
	case !p.context.renderMD:
		return true, nil
	case !allAreText:
		return true, nil
	case !p.context.acceptsHTML:
		return true, nil
	}

	var output []string
	var titles []string

	authors := map[string]struct{}{}

	for _, i := range p.rwqs {
		rwq := i

		output = append(output, string(markdown.ToHTML(p.contents[rwq.r.ID].bytes, nil, nil)))

		authors[rwq.r.Owner] = struct{}{}

		switch {
		case rwq.r.Name != nil:
			titles = append(titles, *rwq.r.Name)
		default:
			titles = append(titles, rwq.qi.Original)
		}
	}

	display := "inline"
	mimetype := "text/html; charset=utf-8"

	if p.context.forceDownload {
		display = "attachment"
		mimetype = "application/octet-stream"
	}

	display = fmt.Sprintf("%s; filename=%s", display, "markdown.html")

	p.w.Header().Set("Content-Type", mimetype)
	p.w.Header().Set("Content-Disposition", display)

	tpl, err := tt.ParseFiles("resources/markdown.html")
	if err != nil {
		return false, err
	}

	var authorS []string

	for a := range authors {
		authorS = append(authorS, a)
	}

	err = tpl.Execute(p.w, map[string]interface{}{
		"title":   strings.Join(titles, ", "),
		"author":  strings.Join(authorS, ", "),
		"content": strings.Join(output, "<br />"),
	})

	return false, err
}
