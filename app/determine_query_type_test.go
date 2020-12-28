package app

import (
	"strings"
	"testing"
)

func TestExtensionsMatch(t *testing.T) {
	suite := []struct {
		Name  string
		Exts  Extensions
		Input string
		Match bool
	}{
		{
			Name:  "matches single extension",
			Exts:  Extensions([]string{"md"}),
			Input: "md",
			Match: true,
		},
		{
			Name:  "matches multiple extensions",
			Exts:  Extensions([]string{"md", "html"}),
			Input: "md.html",
			Match: true,
		},
		{
			Name:  "does not match too many",
			Exts:  Extensions([]string{"md", "html"}),
			Input: "php.md.html",
			Match: false,
		},
		{
			Name:  "does match too few",
			Exts:  Extensions([]string{"php", "md", "html"}),
			Input: "html",
			Match: true,
		},
		{
			Name:  "does not match",
			Exts:  Extensions([]string{"php", "md", "html"}),
			Input: "raw.go",
			Match: false,
		},
	}

	for _, test := range suite {
		i := strings.Split(test.Input, ".")

		if test.Exts.Match(i...) != test.Match {
			t.Errorf("unexpected extensions result for test \"%s\"", test.Name)
		}
	}
}

func TestQIFileName(t *testing.T) {
	a := &App{}

	suite := []struct {
		Name          string
		Input         string
		ExpectedCount int
		OutcomeSet    []string
	}{
		{
			Name:          "named files",
			Input:         "@name.txt",
			ExpectedCount: 1,
			OutcomeSet:    []string{"name.txt"},
		},
		{
			Name:          "hash with extension",
			Input:         "a3f.md",
			ExpectedCount: 1,
			OutcomeSet:    []string{"a3f.md"},
		},
		{
			Name:          "hash without extension",
			Input:         "d2g",
			ExpectedCount: 1,
			OutcomeSet:    []string{"d2g"},
		},
		{
			Name:          "multi query hashes",
			Input:         "af3,gM2,a66",
			ExpectedCount: 3,
			OutcomeSet:    []string{"af3", "gM2", "a66"},
		},
		{
			Name:          "multi query hashes, one with ext",
			Input:         "af3,gM2.json,a66",
			ExpectedCount: 3,
			OutcomeSet:    []string{"af3", "gM2.json", "a66"},
		},
		{
			Name:          "multi query hashes, one with many ext",
			Input:         "af3,gM2.md.html,a66",
			ExpectedCount: 3,
			OutcomeSet:    []string{"af3", "gM2.md.html", "a66"},
		},
		{
			Name:          "multi type, varied",
			Input:         ":alva,a2T.md,@kyle.md,:summin.txt,aaB",
			ExpectedCount: 5,
			OutcomeSet:    []string{"alva", "a2T.md", "kyle.md", "summin.txt", "aaB"},
		},
	}

	for _, test := range suite {
		qis := a.ParseQueryType(test.Input)
		if len(qis) != test.ExpectedCount {
			t.Errorf("failed test %s, wrong query length", test.Name)
		}

		for i, qi := range qis {
			if qi.Filename() != test.OutcomeSet[i] {
				t.Errorf("failed test %s, wrong filename: expecting %s, got %s", test.Name, test.OutcomeSet[i], qi.Filename())
			}
		}
	}
}

func TestQueryInput(t *testing.T) {
	a := &App{}

	suite := []struct {
		Name          string
		Input         string
		ExpectedCount int
		OutcomeSet    []*QueryInput
	}{
		{
			Name:          "single hash",
			Input:         "AAb",
			ExpectedCount: 1,
			OutcomeSet: []*QueryInput{
				{
					QueryType: Hash,
					Original:  "AAb",
					Input:     "AAb",
					Exts:      Extensions([]string{}),
				},
			},
		},
		{
			Name:          "multi hash",
			Input:         "AAb,afG,bbG",
			ExpectedCount: 3,
			OutcomeSet: []*QueryInput{
				{
					QueryType: Hash,
					Original:  "AAb",
					Input:     "AAb",
					Exts:      Extensions([]string{}),
				},
				{
					QueryType: Hash,
					Original:  "afG",
					Input:     "afG",
					Exts:      Extensions([]string{}),
				},
				{
					QueryType: Hash,
					Original:  "bbG",
					Input:     "bbG",
					Exts:      Extensions([]string{}),
				},
			},
		},
		{
			Name:          "multi varient",
			Input:         "@alva.png,:kyle,a2F.md.html,@duffleman,:test.go,A2O",
			ExpectedCount: 6,
			OutcomeSet: []*QueryInput{
				{
					QueryType: Name,
					Original:  "@alva.png",
					Input:     "alva",
					Exts:      Extensions([]string{"png"}),
				},
				{
					QueryType: Shortcut,
					Original:  ":kyle",
					Input:     "kyle",
					Exts:      Extensions([]string{}),
				},
				{
					QueryType: Hash,
					Original:  "a2F.md.html",
					Input:     "a2F",
					Exts:      Extensions([]string{"md", "html"}),
				},
				{
					QueryType: Name,
					Original:  "@duffleman",
					Input:     "duffleman",
					Exts:      Extensions([]string{}),
				},
				{
					QueryType: Shortcut,
					Original:  ":test.go",
					Input:     "test",
					Exts:      Extensions([]string{"go"}),
				},
				{
					QueryType: Hash,
					Original:  "A2O",
					Input:     "A2O",
					Exts:      Extensions([]string{}),
				},
			},
		},
	}

	for _, test := range suite {
		qis := a.ParseQueryType(test.Input)

		if len(qis) != test.ExpectedCount {
			t.Errorf("failed test %s, wrong query length", test.Name)
		}

		for i, qi := range qis {
			matchQI := test.OutcomeSet[i]

			if qi.QueryType != matchQI.QueryType {
				t.Errorf("failed test %s:%d, wrong query type: got %d, expected %d", test.Name, i, qi.QueryType, matchQI.QueryType)
			}

			if qi.Original != matchQI.Original {
				t.Errorf("failed test %s:%d, wrong original: got %s, expected %s", test.Name, i, qi.Original, matchQI.Original)
			}

			if qi.Input != matchQI.Input {
				t.Errorf("failed test %s:%d, wrong input: got %s, expected %s", test.Name, i, qi.Input, matchQI.Input)
			}

			if len(matchQI.Exts) != len(qi.Exts) {
				t.Errorf("failed test %s:%d, extension count wrong: got %d, expected %d", test.Name, i, len(qi.Exts), len(matchQI.Exts))
			}

			for j, ext := range qi.Exts {
				matchExt := matchQI.Exts[j]

				if matchExt != ext {
					t.Errorf("failed test %s:%d:%d, wrong ext: got %s, expected %s", test.Name, i, j, ext, matchExt)
				}
			}
		}
	}
}
