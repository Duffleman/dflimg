package pipeline

import (
	"html/template"

	"dflimg/lib/cher"

	"golang.org/x/sync/errgroup"
)

// LoadFilesFromFS loads files from the filesystem into memory
func LoadFilesFromFS(p *Pipeline) (bool, error) {
	g, gctx := errgroup.WithContext(p.ctx)

	for _, i := range p.rwqs {
		rwq := i

		g.Go(func() (err error) {
			b, modtime, err := p.app.GetFile(gctx, rwq.r)
			if err != nil {
				return err
			}

			p.Lock()
			defer p.Unlock()

			p.contents[rwq.r.ID] = fileContent{
				modtime: modtime,
				bytes:   b,
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
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

	return true, nil
}
