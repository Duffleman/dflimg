package dflimg

import "time"

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
	PostgresCS = "postgres://duffleman@localhost:5432/dflimg?sslmode=disable"
	// DefaultAddr is the default address to listen on
	DefaultAddr = ":3000"
)

// File is a file entity to host in S3
type File struct {
	ID        string    `json:"id"`
	Serial    int       `json:"serial"`
	Owner     string    `json:"owner"`
	S3        string    `json:"s3"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Shortcuts []string  `json:"shortcuts"`
}

// Link is a URL shortener entity
type Link struct {
	ID        string    `json:"id"`
	Owner     string    `json:"owner"`
	URL       string    `json:"url"`
	NSFW      bool      `json:"nsfw"`
	Shortcuts []string  `json:"shortcuts"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// UploadFileResponse is a response for the file upload endpoint
type UploadFileResponse struct {
	FileID string `json:"file_id"`
	Hash   string `json:"hash"`
	URL    string `json:"url"`
}
