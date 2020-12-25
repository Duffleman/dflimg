package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"dflimg"
	dflapp "dflimg/app"
	"dflimg/app/storageproviders"
	dfldb "dflimg/db"
	dflrpc "dflimg/rpc"
	dflmw "dflimg/rpc/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/cors"
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
	// storage provider
	storageProvider := os.Getenv("STORAGE_PROVIDER")

	var err error
	var sp storageproviders.StorageProvider

	switch storageProvider {
	case "aws":
		sp, err = storageproviders.NewAWSProviderFromEnv()
		if err != nil {
			logger.Fatal(err)
		}
	case "lfs":
		sp, err = storageproviders.NewLFSProviderFromEnv()
		if err != nil {
			logger.Fatal(err)
		}
	default:
		logger.Fatal(errors.New("unsupported_provider"))
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
	redisAddr := dflimg.GetEnv("redis_uri")
	redisOpts, err := redis.ParseURL(redisAddr)
	if err != nil {
		logger.Fatal(err)
	}

	redisClient := redis.NewClient(redisOpts)
	_, err = redisClient.Ping().Result()
	if err != nil {
		logger.Fatal(err)
	}

	redis := dflapp.NewCache(redisClient)

	// Setup app & rpc
	router := chi.NewRouter()
	app := dflapp.New(db, sp, hasher, redis)
	rpc := dflrpc.New(logger, router, app)

	// Add middleware
	rpc.Use(middleware.RequestID)
	rpc.Use(middleware.RealIP)
	rpc.Use(middleware.Recoverer)
	rpc.Use(dflmw.AuthMiddleware(dflimg.GetUsers()))
	rpc.Use(middleware.Timeout(60 * time.Second))
	rpc.Use(func(next http.Handler) http.Handler {
		return cors.AllowAll().Handler(next)
	})

	// define routes
	rpc.Get("/", rpc.Homepage)
	rpc.Get("/favicon.ico", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "Not Found", 404)
		return
	})
	rpc.Get("/health", rpc.HealthCheck)
	rpc.Get("/robots.txt", rpc.Robots)

	rpc.Post("/add_shortcut", rpc.AddShortcut)
	rpc.Post("/create_signed_url", rpc.CreateSignedURL)
	rpc.Post("/delete_resource", rpc.DeleteResource)
	rpc.Post("/list_resources", rpc.ListResources)
	rpc.Post("/remove_shortcut", rpc.RemoveShortcut)
	rpc.Post("/resave_hashes", rpc.ResaveHashes)
	rpc.Post("/set_nsfw", rpc.SetNSFW)
	rpc.Post("/shorten_url", rpc.ShortenURL)
	rpc.Post("/upload_file", rpc.UploadFile)
	rpc.Post("/view_details", rpc.ViewDetails)

	rpc.Head(fmt.Sprintf("/%s{query}", dflapp.NameCharacter), rpc.HeadResource)
	rpc.Head("/{query}", rpc.HeadResource)
	rpc.Get(fmt.Sprintf("/%s{query}", dflapp.NameCharacter), rpc.HandleResource)
	rpc.Get("/{query}", rpc.HandleResource)

	// serve
	addr := dflimg.GetEnv("addr")
	rpc.Serve(addr)
}
