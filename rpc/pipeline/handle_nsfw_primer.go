package pipeline

import (
	"html/template"
)

func HandleNSFWPrimer(p *Pipeline) (bool, error) {
	// skip this step if qualifiers are met
	if p.context.forceDownload || p.context.primed || p.context.multifile || !p.rwqs[0].r.NSFW {
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
