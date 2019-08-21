package app

import (
	"bytes"
	"context"
	"dflimg/dflerr"

	"github.com/go-pg/pg"
)

// GetFileByLabel gets a file by it's label
func (a *App) GetFileByShortcut(ctx context.Context, shortcut string) (string, *bytes.Buffer, error) {
	file, err := a.db.FindFileByShortcut(ctx, shortcut)
	if err != nil {
		if err == pg.ErrNoRows {
			return "", nil, dflerr.New(dflerr.NotFound, nil)
		}
		return "", nil, err
	}

	return a.getFileBySerial(ctx, file.Serial)
}
