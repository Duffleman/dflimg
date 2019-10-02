package app

import (
	"dflimg/db"

	"github.com/aws/aws-sdk-go/aws/session"
	hashids "github.com/speps/go-hashids"
)

// App is a struct for the App and it's handlers
type App struct {
	db     *db.DB
	aws    *session.Session
	hasher *hashids.HashID
	redis  *Cache
}

// New returns an instance of the App
func New(db *db.DB, aws *session.Session, hasher *hashids.HashID, redis *Cache) *App {
	return &App{
		db:     db,
		aws:    aws,
		hasher: hasher,
		redis:  redis,
	}
}
