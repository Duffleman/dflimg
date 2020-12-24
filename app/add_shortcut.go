package app

import (
	"context"

	"dflimg"
	"dflimg/db"
	"dflimg/lib/cher"
)

func (a *App) AddShortcut(ctx context.Context, resource *dflimg.Resource, shortcut string) error {
	err := a.db.FindShortcutConflicts(ctx, []string{shortcut})
	if err != nil {
		return cher.New("shortcuts_already_taken", cher.M{"shortcut": shortcut}, cher.Coerce(err))
	}

	return a.db.ChangeShortcut(ctx, db.ArrayAdd, resource.ID, shortcut)
}
