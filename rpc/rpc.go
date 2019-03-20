package rpc

import (
	"encoding/json"
	"errors"
	"net/http"

	"dflimg"
	"dflimg/app"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

var (
	// ErrAccessDenied is an error to show that access is denied
	ErrAccessDenied = errors.New("access denied")
)

// RPC is a struct for the RPC server and it's handlers
type RPC struct {
	logger *logrus.Logger
	router *chi.Mux
	app    *app.App
}

// New returns a new instance of an RPC handler
func New(logger *logrus.Logger, router *chi.Mux, app *app.App) *RPC {
	return &RPC{
		logger: logger,
		app:    app,
		router: router,
	}
}

// Use is a wrapper for chi's Use func
func (r *RPC) Use(middlewares ...func(http.Handler) http.Handler) {
	r.router.Use(middlewares...)
}

// Get is a wrapper for chi's Get func
func (r *RPC) Get(pattern string, handlerFn http.HandlerFunc) {
	r.router.Get(pattern, handlerFn)
}

// Post is a wrapper for chi's Post func
func (r *RPC) Post(pattern string, handlerFn http.HandlerFunc) {
	r.router.Post(pattern, handlerFn)
}

// Serve starts the HTTP server
func (r *RPC) Serve(port string) {
	r.logger.Info("starting web server")
	http.ListenAndServe(port, r.router)
}

func (r *RPC) handleError(w http.ResponseWriter, req *http.Request, err error, meta *map[string]interface{}) {
	l := logrus.NewEntry(r.logger)

	if meta != nil {
		l = l.WithFields(logrus.Fields(*meta))
	}

	switch err {
	case dflimg.ErrNotFound:
		l.Info(err)
		w.WriteHeader(404)
	case ErrAccessDenied:
		l.Info(err)
		w.WriteHeader(403)
	default:
		l.Warn(err)
		w.WriteHeader(500)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": err.Error(),
		"meta": err,
	})

	return
}
