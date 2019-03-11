package main

import (
	"time"

	"dflimg/app"
	dfldb "dflimg/db"
	dflrpc "dflimg/rpc"
	dflmw "dflimg/rpc/middleware"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
)

func main() {
	// Setup logger
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
	}

	// user config for username vs API key
	users := map[string]string{
		"Duffleman": "test",
	}

	// setup app dependancies
	// aws
	s, err := session.NewSession(&aws.Config{Region: aws.String(app.S3Region)})
	if err != nil {
		logger.Fatal(err)
	}

	// postgres db
	pgdb := pg.Connect(&pg.Options{
		User:     "duffleman",
		Database: "dflimg",
	})
	defer pgdb.Close()

	db := dfldb.New(pgdb)

	// Setup app & rpc
	router := chi.NewRouter()
	app := app.New(db, s)
	rpc := dflrpc.New(logger, router, app)

	// Add middleware
	rpc.Use(middleware.RequestID)
	rpc.Use(middleware.RealIP)
	rpc.Use(middleware.Recoverer)
	rpc.Use(dflmw.AuthMiddleware(users))
	rpc.Use(middleware.Timeout(60 * time.Second))

	// define routes
	rpc.Get("/health", rpc.HealthCheck)
	rpc.Post("/upload", rpc.Upload)

	// serve
	rpc.Serve(":3000")
}
