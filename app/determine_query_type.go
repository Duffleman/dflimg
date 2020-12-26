package app

import (
	"fmt"
	"strings"
)

// ShortcutCharacter marks the character used to find shortcuts
const ShortcutCharacter = ":"

// NameCharacter marks the character used to find files by their names
const NameCharacter = "@"

type QueryType int

const (
	Hash QueryType = iota
	Shortcut
	Name
)

type Extensions []string

func (e Extensions) Match(in ...string) bool {
	if len(in) > len(e) {
		return false
	}

	matchSet := strings.Join(e[len(e)-len(in):], ".")

	return matchSet == strings.Join(in, ".")
}

func (e Extensions) Last() string {
	if len(e) == 0 {
		return ""
	}

	return e[len(e)-1]
}

type QueryInput struct {
	QueryType QueryType
	Original  string
	Input     string
	Exts      Extensions
}

func (qi QueryInput) Filename() string {
	if len(qi.Exts) == 0 {
		return qi.Input
	}

	extensions := strings.Join(qi.Exts, ".")

	return fmt.Sprintf("%s.%s", qi.Input, extensions)
}

func (a *App) ParseQueryType(inStr string) (q []*QueryInput) {
	for _, query := range strings.Split(inStr, ",") {
		var exts []string
		var qt QueryType = Hash
		var input = query

		if strings.ContainsRune(query, '.') {
			parts := strings.Split(query, ".")

			input = parts[0]
			exts = parts[1:]
		}

		if strings.HasPrefix(query, NameCharacter) {
			qt = Name
			input = strings.TrimPrefix(input, NameCharacter)
		}

		if strings.HasPrefix(query, ShortcutCharacter) {
			qt = Shortcut
			input = strings.TrimPrefix(input, ShortcutCharacter)
		}

		q = append(q, &QueryInput{
			QueryType: qt,
			Original:  query,
			Input:     input,
			Exts:      Extensions(exts),
		})
	}

	return
}
