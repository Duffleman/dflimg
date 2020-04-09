package main

import (
	"context"
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
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v4/pgxpool"
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

	// database (postgres)
	poolConfig, err := pgxpool.ParseConfig(dflimg.GetEnv("pg_connection_string"))
	if err != nil {
		logger.Fatal(err)
	}

	pgdb, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Fatal(err)
	}

	db := dfldb.New(pgdb)

	// hasher
	hd := hashids.NewData()

	salt := dflimg.GetEnv("salt")

	hd.Salt = salt
	hd.MinLength = dflimg.EncodeLength

	hasher, _ := hashids.NewWithData(hd)

	// Cache
	// cache := cache.New(30*time.Minute, 1*time.Hour)
	redisAddr := dflimg.GetEnv("redis_addr")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	_, err = redisClient.Ping().Result()
	if err != nil {
		logger.Fatal(err)
	}
	redis := app.NewCache(redisClient)

	// Setup app & rpc
	router := chi.NewRouter()
	app := app.New(db, s, hasher, redis)
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
	rpc.Get("/robots.txt", rpc.Robots)
	rpc.Post("/upload_file", rpc.UploadFile)
	rpc.Post("/create_signed_url", rpc.CreateSignedURL)
	rpc.Post("/shorten_url", rpc.ShortenURL)
	rpc.Post("/delete_resource", rpc.DeleteResource)
	rpc.Post("/resave_hashes", rpc.ResaveHashes)
	rpc.Get("/{query}", rpc.GetResource)

	// serve
	addr := dflimg.GetEnv("addr")
	rpc.Serve(addr)
}
