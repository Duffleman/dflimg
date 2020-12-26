package pipeline

import (
	"dflimg/app"

	"html/template"
)

func HandleSyntaxHighlighting(p *Pipeline) (bool, error) {
	switch {
	case p.context.multifile:
		return true, nil
	case p.context.forceDownload:
		return true, nil
	case !p.rwqs[0].context.isText:
		return true, nil
	case !p.context.acceptsHTML:
		return true, nil
	case p.rwqs[0].qi.QueryType != app.Name && p.rwqs[0].qi.Ext == nil:
		return true, nil
	case p.rwqs[0].qi.QueryType == app.Name && !p.context.wantsHighlighting:
		return true, nil
	}

	tpl, err := template.ParseFiles("resources/code.html")
	if err != nil {
		return false, err
	}

	var language string

	if p.rwqs[0].qi.Ext != nil {
		language = *p.rwqs[0].qi.Ext
	}

	if p.context.highlightLanguage != "" {
		language = p.context.highlightLanguage
	}

	var name string = *p.rwqs[0].r.Hash

	if p.rwqs[0].r.Name != nil {
		name = *p.rwqs[0].r.Name
	}

	err = tpl.Execute(p.w, map[string]interface{}{
		"language": language,
		"title":    name,
		"author":   p.rwqs[0].r.Owner,
		"content":  string(p.contents[p.rwqs[0].r.ID].bytes),
	})
	return false, err
}
