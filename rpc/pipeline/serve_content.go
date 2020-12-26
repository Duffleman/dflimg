package pipeline

import (
	"bytes"
	"fmt"
	"net/http"

	"dflimg/lib/cher"
)

func ServeContent(p *Pipeline) (bool, error) {
	if p.context.multifile {
		return false, cher.New("multiple_files", nil)
	}

	rwq := p.rwqs[0]

	display := "inline"
	mimetype := "application/octet-stream"

	if p.context.forceDownload {
		display = "attachment"
	} else if rwq.r.MimeType != nil {
		mimetype = *rwq.r.MimeType
	}

	if rwq.r.Name != nil {
		display = fmt.Sprintf("%s; filename=%s", display, *rwq.r.Name)
	}

	p.w.Header().Set("Content-Type", mimetype)
	p.w.Header().Set("Content-Disposition", display)

	reader := bytes.NewReader(p.contents[rwq.r.ID].bytes)

	http.ServeContent(p.w, p.r, rwq.qi.Filename(), *p.contents[rwq.r.ID].modtime, reader)

	return true, nil
}
