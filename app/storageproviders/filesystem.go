package storageproviders

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"dflimg"

	"golang.org/x/sync/errgroup"
)

// LocalFileSystem is a storage provider for the local filesystem
type LocalFileSystem struct {
	folder      string
	permissions os.FileMode
}

// NewLFSProviderFromEnv makes a new FileSystem provider from env vars
func NewLFSProviderFromEnv() (StorageProvider, error) {
	folder := os.Getenv("LFS_FOLDER")
	if folder == "" {
		return nil, errors.New("missing_lfs_folder")
	}

	var permissions os.FileMode = 0777

	permissionsStr := os.Getenv("LFS_PERMISSIONS")
	if permissionsStr != "" {
		i, err := strconv.Atoi(permissionsStr)
		if err != nil {
			return nil, err
		}

		permissions = os.FileMode(i)
	}

	return &LocalFileSystem{
		folder:      folder,
		permissions: permissions,
	}, nil
}

// GenerateKey generates a file key used as a filename
func (fs *LocalFileSystem) GenerateKey(fileID string) string {
	return fmt.Sprintf("%s/%s", fs.folder, fileID)
}

// SupportsSignedURLs lets the service know if it can use prepared URLs
func (fs *LocalFileSystem) SupportsSignedURLs() bool {
	return false
}

// Get a resource from the storage provider
func (fs *LocalFileSystem) Get(ctx context.Context, resource *dflimg.Resource) ([]byte, *time.Time, error) {
	g, _ := errgroup.WithContext(ctx)

	var bytes []byte
	var lastModified time.Time

	g.Go(func() (err error) {
		var fileInfo os.FileInfo

		fileInfo, err = os.Stat(resource.Link)
		if err != nil {
			return err
		}

		lastModified = fileInfo.ModTime()

		return
	})

	g.Go(func() (err error) {
		bytes, err = ioutil.ReadFile(resource.Link)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, nil, err
	}

	return bytes, &lastModified, nil
}

// PrepareUpload prepares an upload into the storage provider
func (fs *LocalFileSystem) PrepareUpload(ctx context.Context, key, contentType string, expiry time.Duration) (string, error) {
	return "", errors.New("unsupported")
}

// Upload a resource into the storage provider
func (fs *LocalFileSystem) Upload(_ context.Context, key, contentType string, file bytes.Buffer) error {
	return ioutil.WriteFile(key, file.Bytes(), fs.permissions)
}
