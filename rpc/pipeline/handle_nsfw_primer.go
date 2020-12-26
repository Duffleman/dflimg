package pipeline

import (
	"html/template"
)

// HandleNSFWPrimer will show a NSFW primer screen if a resource has NSFW content
func HandleNSFWPrimer(p *Pipeline) (bool, error) {
	if p.context.forceDownload || p.context.primed {
		return true, nil
	}

	var anyNSFW bool

	for _, i := range p.rwqs {
		rwq := i

		if rwq.r.NSFW {
			anyNSFW = true
		}
	}

	if !anyNSFW {
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
