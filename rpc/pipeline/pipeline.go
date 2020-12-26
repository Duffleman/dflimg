package pipeline

import (
	"context"
	"net/http"
	"time"

	"dflimg"
	"dflimg/app"
)

type HandlerType func(*Pipeline) (bool, error)

type resourceWithQuery struct {
	r       *dflimg.Resource
	qi      *app.QueryInput
	context struct {
		isText bool
	}
}

type fileContent struct {
	modtime *time.Time
	bytes   []byte
}

type Pipeline struct {
	ctx      context.Context
	app      *app.App
	r        *http.Request
	w        http.ResponseWriter
	qi       []*app.QueryInput
	rwqs     []*resourceWithQuery
	contents map[string]fileContent
	context  struct {
		forceDownload     bool
		primed            bool
		wantsHighlighting bool
		highlightLanguage string
		multifile         bool
		acceptsHTML       bool
	}
	steps []HandlerType
}

type Creator struct {
	App *app.App
	QI  []*app.QueryInput
	R   *http.Request
	W   http.ResponseWriter
}

func New(ctx context.Context, pc *Creator) *Pipeline {
	return &Pipeline{
		ctx:      ctx,
		qi:       pc.QI,
		app:      pc.App,
		r:        pc.R,
		w:        pc.W,
		contents: make(map[string]fileContent),
	}
}

func (p *Pipeline) Steps(steps []HandlerType) {
	p.steps = steps
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
