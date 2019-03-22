package app

import (
	"bytes"
	"context"
	"dflimg/dflerr"

	"github.com/go-pg/pg"
)

// GetFileByLabel gets a file by it's label
func (a *App) GetFileByLabel(ctx context.Context, label string) (string, *bytes.Buffer, error) {
	file, err := a.db.FindFileByLabel(label)
	if err != nil {
		if err == pg.ErrNoRows {
			return "", nil, dflerr.New(dflerr.NotFound, nil)
		}
		return "", nil, err
	}

	return a.getFileBySerial(ctx, file.Serial)
}
