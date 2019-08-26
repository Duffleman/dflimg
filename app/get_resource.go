package app

import (
	"context"
	"strings"

	"dflimg"
)

// ShortcutCharacter marks the character used to find shortcuts
const ShortcutCharacter = ":"

func (a *App) GetResource(ctx context.Context, input string) (*dflimg.Resource, error) {
	if strings.HasPrefix(input, ShortcutCharacter) {
		return a.GetResourceByShortcut(ctx, input)
	} else {
		return a.GetResourceByHash(ctx, input)
	}
}
