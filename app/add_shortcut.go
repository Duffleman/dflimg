package app

import (
	"context"

	"dflimg"
	"dflimg/db"
	"dflimg/dflerr"
)

func (a *App) AddShortcut(ctx context.Context, resource *dflimg.Resource, shortcut string) error {
	err := a.db.FindShortcutConflicts(ctx, []string{shortcut})
	if err != nil {
		return dflerr.New("shortcuts_already_taken", dflerr.M{"shortcut": shortcut}, dflerr.Parse(err))
	}

	return a.db.ChangeShortcut(ctx, db.ArrayAdd, resource.ID, shortcut)
}
