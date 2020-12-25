package app

import (
	"fmt"
	"strings"
)

type QueryType int

const (
	Hash QueryType = iota
	Shortcut
	Name
)

type QueryInput struct {
	QueryType QueryType
	Original  string
	Input     string
	Ext       *string
}

func (qi QueryInput) Filename() string {
	if qi.Ext == nil {
		return qi.Input
	}

	return fmt.Sprintf("%s.%s", qi.Input, *qi.Ext)
}

// ShortcutCharacter marks the character used to find shortcuts
const ShortcutCharacter = ":"

// NameCharacter marks the character used to find files by their names, must have a trailing slash
const NameCharacter = "n/"

func (a *App) ParseQueryType(query string) *QueryInput {
	var ext *string
	var qt QueryType = Hash
	var input = query

	if strings.ContainsRune(query, '.') {
		parts := strings.Split(query, ".")
		ext = &parts[len(parts)-1]
		input = strings.Join(parts[:len(parts)-1], ".")
	}

	if strings.HasPrefix(query, NameCharacter) {
		qt = Name
		input = strings.TrimPrefix(input, NameCharacter)
	}

	if strings.HasPrefix(query, ShortcutCharacter) {
		qt = Shortcut
		input = strings.TrimPrefix(input, ShortcutCharacter)
	}

	return &QueryInput{
		QueryType: qt,
		Original:  query,
		Input:     input,
		Ext:       ext,
	}
}
