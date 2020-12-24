package app

import (
	"context"
	"strings"

	"dflimg"
)

// ShortcutCharacter marks the character used to find shortcuts
const ShortcutCharacter = ":"

// NameCharacter marks the character used to find files by their names, must have a trailing slash
const NameCharacter = "n/"

// GetResource returns a resource by shortcut or hash. Regardless of it's deleted status
func (a *App) GetResource(ctx context.Context, input string) (res *dflimg.Resource, ext *string, err error) {
	rootURL := dflimg.GetEnv("root_url") + "/"

	if strings.HasPrefix(input, rootURL) {
		input = strings.TrimPrefix(input, rootURL)
	}

	if strings.HasPrefix(input, NameCharacter) {
		input = strings.TrimPrefix(input, NameCharacter)
		res, err = a.GetResourceByName(ctx, input)
		return
	}

	if strings.ContainsRune(input, '.') {
		parts := strings.Split(input, ".")

		ext = &parts[len(parts)-1]
		input = strings.Join(parts[:len(parts)-1], ".")
	}

	if strings.HasPrefix(input, ShortcutCharacter) {
		res, err = a.GetResourceByShortcut(ctx, input)
		return
	}

	res, err = a.GetResourceByHash(ctx, input)
	return
}
