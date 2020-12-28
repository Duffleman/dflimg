package pipeline

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"dflimg"
	"dflimg/app"
)

func TestParsesContext(t *testing.T) {
	ctx := context.Background()

	rootURL := "https://dfl.mn"
	app := &app.App{}

	suite := []struct {
		Name         string
		Input        string
		MatchContext pipelineContext
		Headers      http.Header
	}{
		{
			Name:  "simple hash",
			Input: "aab",
		},
		{
			Name:         "will force download with ?d",
			Input:        "@myfile.png?d",
			MatchContext: pipelineContext{forceDownload: true},
		},
		{
			Name:         "is primed with ?primed",
			Input:        "@myfile.png?primed",
			MatchContext: pipelineContext{primed: true},
		},
		{
			Name:         "forces highlighting with ?sh if accepts html",
			Input:        "aS2?sh",
			MatchContext: pipelineContext{wantsHighlighting: true},
			Headers:      http.Header{"Accept": []string{"text/html"}},
		},
		{
			Name:         "forces highlighting with ?sh and specific language if accepts html",
			Input:        "aS2?sh=json",
			MatchContext: pipelineContext{wantsHighlighting: true, highlightLanguage: "json"},
			Headers:      http.Header{"Accept": []string{"text/html"}},
		},
		{
			Name:         "will not highlight with ?sh and specific language without accepting html",
			Input:        "aS2?sh",
			MatchContext: pipelineContext{wantsHighlighting: false},
		},
		{
			Name:         "understands multifile",
			Input:        "aMb,@alva,@kyle,:duffleman,m99.json",
			MatchContext: pipelineContext{multifile: true},
		},
		{
			Name:         "understands single files with tricks",
			Input:        "aMb,",
			MatchContext: pipelineContext{multifile: false},
		},
		{
			Name:         "will render MD for single files with .md.html if accepts html",
			Input:        "@hello.md.html",
			MatchContext: pipelineContext{renderMD: true},
			Headers:      http.Header{"Accept": []string{"text/html"}},
		},
		{
			Name:         "will not render MD for single files with .md.html without accepting html",
			Input:        "@hello.md.html",
			MatchContext: pipelineContext{renderMD: false},
		},
		{
			Name:         "will render MD for multi files with ?pmd if accepts html",
			Input:        "afF,abS,@alva.md?pmd",
			MatchContext: pipelineContext{renderMD: true, multifile: true},
			Headers:      http.Header{"Accept": []string{"text/html"}},
		},
		{
			Name:         "will not render MD for multi files with ?pmd without accepting html",
			Input:        "afF,abS,@alva.md?pmd",
			MatchContext: pipelineContext{renderMD: false, multifile: true},
		},
		{
			Name:         "render MD on single file",
			Input:        "aab.md.html",
			MatchContext: pipelineContext{renderMD: true, forceDownload: false},
			Headers:      http.Header{"Accept": []string{"text/html"}},
		},
		{
			Name:         "render MD on multi file",
			Input:        "aab,a2M,@alva?pmd",
			MatchContext: pipelineContext{renderMD: true, forceDownload: false, multifile: true},
			Headers:      http.Header{"Accept": []string{"text/html"}},
		},
	}

	for _, test := range suite {
		request, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", rootURL, test.Input), nil)
		if err != nil {
			t.Errorf("cannot make request: %w", err)
		}

		request.Header = test.Headers

		pipe := New(ctx, &Creator{
			App: app,
			R:   request,
		})

		qis := pipe.app.ParseQueryType(test.Input)

		for _, qi := range qis {
			pipe.rwqs = append(pipe.rwqs, &resourceWithQuery{
				r:  &dflimg.Resource{},
				qi: qi,
			})
		}

		pipe.Steps([]HandlerType{
			MakeContext,
		})

		err = pipe.Run()
		if err != nil {
			t.Errorf("failed to run pipe; %w", err)
		}

		if pipe.context.forceDownload != test.MatchContext.forceDownload {
			t.Errorf("failed test %s, wrong context for forceDownload: got %t, expected %t", test.Name, pipe.context.forceDownload, test.MatchContext.forceDownload)
		}

		if pipe.context.primed != test.MatchContext.primed {
			t.Errorf("failed test %s, wrong context for primed: got %t, expected %t", test.Name, pipe.context.primed, test.MatchContext.primed)
		}

		if pipe.context.wantsHighlighting != test.MatchContext.wantsHighlighting {
			t.Errorf("failed test %s, wrong context for wantsHighlighting: got %t, expected %t", test.Name, pipe.context.wantsHighlighting, test.MatchContext.wantsHighlighting)
		}

		if pipe.context.highlightLanguage != test.MatchContext.highlightLanguage {
			t.Errorf("failed test %s, wrong context for highlightLanguage: got %s, expected %s", test.Name, pipe.context.highlightLanguage, test.MatchContext.highlightLanguage)
		}

		if pipe.context.multifile != test.MatchContext.multifile {
			t.Errorf("failed test %s, wrong context for multifile: got %t, expected %t", test.Name, pipe.context.multifile, test.MatchContext.multifile)
		}

		if pipe.context.renderMD != test.MatchContext.renderMD {
			t.Errorf("failed test %s, wrong context for renderMD: got %t, expected %t", test.Name, pipe.context.renderMD, test.MatchContext.renderMD)
		}
	}
}
