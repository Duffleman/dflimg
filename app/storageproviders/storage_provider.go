package storageproviders

import (
	"bytes"
	"context"
	"time"

	"dflimg"
)

// StorageProvider is an interface all custom defined storage providers must conform to
type StorageProvider interface {
	CheckEnvVariables() error
	GenerateKey(string) string
	SupportsSignedURLs() bool
	Get(context.Context, *dflimg.Resource) ([]byte, *time.Time, error)
	PrepareUpload(ctx context.Context, key, contentType string, expiry time.Duration) (string, error)
	Upload(ctx context.Context, key, contentType string, file bytes.Buffer) error
}
