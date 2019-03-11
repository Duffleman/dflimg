package app

import (
	"dflimg/db"

	"github.com/aws/aws-sdk-go/aws/session"
)

// App is a struct for the App and it's handlers
type App struct {
	db  *db.DB
	aws *session.Session
}

const (
	S3Region = "eu-west-1"
	S3Bucket = "i.dfl.mn"
)

// New returns an instance of the App
func New(db *db.DB, aws *session.Session) *App {
	return &App{
		db:  db,
		aws: aws,
	}
}
