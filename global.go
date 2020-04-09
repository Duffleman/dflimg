package dflimg

import (
	"bytes"
	"time"
)

// Users is a map[string]string for users to upload keys
var Users = map[string]string{
	"Duffleman": "test",
}

const (
	// S3Region is the region to upload to in S3
	S3Region = "eu-west-1"
	// S3Bucket is the name of the bucket to upload to in S3
	S3Bucket = "s3.duffleman.co.uk"
	// S3RootKey is the folder that stores the images inside the bucket
	S3RootKey = "i.dfl.mn"
	// Salt for encoding the serial IDs to hashes
	Salt = "savour-shingle-sidney-rajah-punk-lead-jenny-scot"
	// EncodeLength - length of the outputted URL (minimum)
	EncodeLength = 3
	// RootURL is the root URL this service runs as
	RootURL = "http://localhost:3000"
	// PostgresCS is the default connection string
	PostgresCS = "postgres://postgres@localhost:5432/dflimg?sslmode=disable"
	// DefaultAddr is the default address to listen on
	DefaultAddr = ":3000"
	// DefaultRedisAddr is the default redis instance to connect to
	DefaultRedisAddr = "localhost:6379"
)

type Resource struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	Serial    int        `json:"serial"`
	Hash      *string    `json:"hash"`
	Owner     string     `json:"owner"`
	Link      string     `json:"link"`
	NSFW      bool       `json:"nsfw"`
	MimeType  *string    `json:"mime_type"`
	Shortcuts []string   `json:"shortcuts"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type ShortFormResource struct {
	ID     string
	Serial int
}

type CreateFileRequest struct {
	File bytes.Buffer `json:"file"`
}

type CreateURLRequest struct {
	URL string `json:"url"`
}

type CreateResourceResponse struct {
	ResourceID string `json:"resource_id"`
	Type       string `json:"type"`
	Hash       string `json:"hash"`
	URL        string `json:"url"`
}

type CreateSignedURLRequest struct {
	ContentType string `json:"content_type"`
}

type CreateSignedURLResponse struct {
	ResourceID string `json:"resource_id"`
	Type       string `json:"type"`
	Hash       string `json:"hash"`
	URL        string `json:"url"`
	S3Link     string `json:"s3link"`
}

type IdentifyResource struct {
	Query string `json:"query"`
}

type SetNSFWRequest struct {
	IdentifyResource
	NSFW bool `json:"nsfw"`
}
