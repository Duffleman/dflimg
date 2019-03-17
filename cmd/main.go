package main

import (
	"time"

	"dflimg"
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
	hashids "github.com/speps/go-hashids"
)

func main() {
	// Setup logger
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
	}

	// setup app dependancies
	// aws
	s, err := session.NewSession(&aws.Config{Region: aws.String(dflimg.S3Region)})
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

	// hasher
	hd := hashids.NewData()

	salt := dflimg.GetEnv("salt")

	hd.Salt = salt
	hd.MinLength = dflimg.EncodeLength

	hasher, _ := hashids.NewWithData(hd)

	// Setup app & rpc
	router := chi.NewRouter()
	app := app.New(db, s, hasher)
	rpc := dflrpc.New(logger, router, app)

	// Add middleware
	rpc.Use(middleware.RequestID)
	rpc.Use(middleware.RealIP)
	rpc.Use(middleware.Recoverer)
	rpc.Use(dflmw.AuthMiddleware(dflimg.Users))
	rpc.Use(middleware.Timeout(60 * time.Second))

	// define routes
	rpc.Get("/", rpc.Homepage)
	rpc.Get("/health", rpc.HealthCheck)
	rpc.Get("/{fileID}", rpc.GetFile)
	rpc.Post("/upload", rpc.Upload)

	// serve
	rpc.Serve(":3000")
}
