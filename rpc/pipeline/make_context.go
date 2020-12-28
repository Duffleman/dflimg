package pipeline

import (
	"strings"
)

// MakeContext generates context info from the request to help future steps in
// the pipeline.
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

	for _, i := range p.rwqs {
		rwq := i

		if rwq.r.MimeType != nil && strings.Contains(*rwq.r.MimeType, "text/plain") {
			rwq.context.isText = true
		}
	}

	if !p.context.multifile && p.rwqs[0].qi.Exts.Match("md", "html") {
		p.context.renderMD = true
	}

	if _, ok := p.r.URL.Query()["pmd"]; ok {
		p.context.renderMD = true
	}

	if !strings.Contains(p.r.Header.Get("Accept"), "text/html") {
		p.context.wantsHighlighting = false
		p.context.renderMD = false
	}

	return true, nil
}
