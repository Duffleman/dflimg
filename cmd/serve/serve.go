package main

import (
	"database/sql"
	"net/http"
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
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	hashids "github.com/speps/go-hashids"
)

func main() {
	// Setup logger
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: false,
	}

	// setup app dependancies
	// aws
	s, err := session.NewSession(&aws.Config{Region: aws.String(dflimg.S3Region)})
	if err != nil {
		logger.Fatal(err)
	}

	pgdb, err := sql.Open("postgres", dflimg.GetEnv("pg_connection_string"))
	if err != nil {
		logger.Fatal(err)
	}
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
	rpc.Use(dflmw.AuthMiddleware(dflimg.GetUsers()))
	rpc.Use(middleware.Timeout(60 * time.Second))

	// define routes
	rpc.Get("/", rpc.Homepage)
	rpc.Get("/favicon.ico", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "Not Found", 404)
		return
	})
	rpc.Get("/health", rpc.HealthCheck)
	rpc.Post("/upload_file", rpc.UploadFile)
	rpc.Post("/shorten_url", rpc.ShortenURL)
	rpc.Get("/list_labels", rpc.ListLabels)
	rpc.Post("/tag_resource", rpc.TagResource)
	rpc.Get("/{input}", rpc.GetResource)

	// serve
	addr := dflimg.GetEnv("addr")
	rpc.Serve(addr)
}
