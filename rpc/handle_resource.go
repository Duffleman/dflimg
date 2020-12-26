package rpc

import (
	"net/http"

	"dflimg/rpc/pipeline"

	"github.com/go-chi/chi"
)

// HandleResource loads a resource and handles the response for it
func (r *RPC) HandleResource(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	qi := r.app.ParseQueryType(chi.URLParam(req, "query"))

	pipe := pipeline.New(ctx, &pipeline.Creator{
		App: r.app,
		R:   req,
		W:   w,
		QI:  qi,
	})

	pipe.Steps([]pipeline.HandlerType{
		// ordered
		pipeline.ListFromDB,
		pipeline.HandleDeleted,
		pipeline.MakeContext,
		pipeline.ValidateRequest,
		pipeline.HandleNSFWPrimer,
		pipeline.HandleURLRedirect,
		pipeline.LoadFilesFromFS,
		pipeline.SyntaxHighlighter,
		pipeline.RenderMD,
		pipeline.FilterMultiFile,
		pipeline.ServeContent,
	})

	r.handleError(w, req, pipe.Run())
}
