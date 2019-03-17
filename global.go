package dflimg

import (
	"errors"
	"os"
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
)

func GetEnv(key string) string {
	var v string
	switch key {
	case "salt":
		v = os.Getenv("DFL_SALT")
		if v == "" {
			return Salt
		}
	case "root_url":
		v = os.Getenv("DFL_ROOT_URL")
		if v == "" {
			return RootURL
		}
	}

	return v
}

var (
	// ErrNotFound is an error for not_found
	ErrNotFound = errors.New("not found")
)

// UploadFileResponse is a response for the file upload endpoint
type UploadFileResponse struct {
	FileID string `json:"file_id"`
	Hash   string `json:"hash"`
	URL    string `json:"url"`
}
