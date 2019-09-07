package app

import (
	"context"
	"strings"

	"dflimg"
)

// ShortcutCharacter marks the character used to find shortcuts
const ShortcutCharacter = ":"

// GetResource returns a resource by shortcut or hash. Regardless of it's deleted status
func (a *App) GetResource(ctx context.Context, input string) (*dflimg.Resource, error) {
	if strings.HasPrefix(input, ShortcutCharacter) {
		return a.GetResourceByShortcut(ctx, input)
	}

	return a.GetResourceByHash(ctx, input)
}
