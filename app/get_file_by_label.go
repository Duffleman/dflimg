package app

import (
	"bytes"
	"context"
	"dflimg"

	"github.com/go-pg/pg"
)

// GetFileByLabel gets a file by it's label
func (a *App) GetFileByLabel(ctx context.Context, label string) (string, *bytes.Buffer, error) {
	file, err := a.db.FindFileByLabel(label)
	if err != nil {
		if err == pg.ErrNoRows {
			return "", nil, dflimg.ErrNotFound
		}
		return "", nil, err
	}

	return a.getFileBySerial(ctx, file.Serial)
}
