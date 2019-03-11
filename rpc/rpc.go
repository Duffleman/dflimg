package rpc

import (
	"dflimg/app"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
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
