package app

import (
	"dflimg/app/storageproviders"
	"dflimg/db"

	hashids "github.com/speps/go-hashids"
)

// App is a struct for the App and it's handlers
type App struct {
	db           *db.DB
	fileProvider storageproviders.StorageProvider
	hasher       *hashids.HashID
	redis        *Cache
}

// New returns an instance of the App
func New(db *db.DB, fileProvider storageproviders.StorageProvider, hasher *hashids.HashID, redis *Cache) *App {
	return &App{
		db:           db,
		fileProvider: fileProvider,
		hasher:       hasher,
		redis:        redis,
	}
}
