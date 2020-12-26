package pipeline

import (
	"strings"
)

func MakeContext(p *Pipeline) (bool, error) {
	if fd := p.r.URL.Query()["d"]; len(fd) >= 1 {
		p.context.forceDownload = true
	}

	if _, ok := p.r.URL.Query()["primed"]; ok {
		p.context.primed = true
	}

	if v, ok := p.r.URL.Query()["sh"]; ok {
		p.context.wantsHighlighting = true
		p.context.highlightLanguage = v[0]
	}

	if len(p.rwqs) > 1 {
		p.context.multifile = true
	}

	for _, rwq := range p.rwqs {
		if rwq.r.MimeType != nil && strings.Contains(*rwq.r.MimeType, "text/plain") {
			rwq.context.isText = true
		}
	}

	if strings.Contains(p.r.Header.Get("Accept"), "text/html") {
		p.context.acceptsHTML = true
	}

	return true, nil
}
