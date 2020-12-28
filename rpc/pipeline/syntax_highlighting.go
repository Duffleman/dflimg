package pipeline

import (
	"html/template"
	"strings"

	"dflimg/app"
)

// SyntaxHighlighter will apply syntax highlighitng to a set of files
func SyntaxHighlighter(p *Pipeline) (bool, error) {
	// if we want to download the files, we won't highlight them
	if p.context.forceDownload {
		return true, nil
	}

	// skip highlighting if we have name queries
	for _, i := range p.rwqs {
		rwq := i

		if rwq.qi.QueryType == app.Name {
			return true, nil
		}
	}

	var atLeastOneExt bool

	// don't highlight files where within the set, one isn't text
	for _, i := range p.rwqs {
		rwq := i

		if !rwq.context.isText {
			return true, nil
		}

		if len(rwq.qi.Exts) >= 1 {
			atLeastOneExt = true
		}
	}

	if !atLeastOneExt {
		return true, nil
	}

	if p.context.renderMD {
		return true, nil
	}

	tpl, err := template.ParseFiles("resources/syntax_highlight.html")
	if err != nil {
		return false, err
	}

	var titles []string
	authorMap := make(map[string]struct{})

	var rs []resourceSet

	for _, i := range p.rwqs {
		rwq := i
		titles = append(titles, rwq.qi.Original)
		authorMap[rwq.r.Owner] = struct{}{}

		language := rwq.qi.Exts.Last()

		if p.context.highlightLanguage != "" {
			language = p.context.highlightLanguage
		}

		rs = append(rs, resourceSet{
			Language: language,
			Content:  string(p.contents[rwq.r.ID].bytes),
		})
	}

	var authors []string

	for a := range authorMap {
		authors = append(authors, a)
	}

	err = tpl.Execute(p.w, map[string]interface{}{
		"resources": rs,
		"title":     strings.Join(titles, ", "),
		"author":    strings.Join(authors, ", "),
	})
	return false, err
}

type resourceSet struct {
	Language string
	Content  string
}
