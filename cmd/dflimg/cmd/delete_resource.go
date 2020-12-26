package cmd

import (
	"context"
	"strings"
	"time"

	"dflimg"
	"dflimg/lib/cher"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var DeleteResourceCmd = &cobra.Command{
	Use:     "delete {query}",
	Aliases: []string{"d"},
	Short:   "Delete a resource",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 || len(args) == 0 {
			return nil
		}

		return cher.New("missing_arguments", nil)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		startTime := time.Now()

		query, err := handleQueryInput(args)
		if err != nil {
			return err
		}

		err = deleteResource(ctx, query)
		if err != nil {
			return err
		}

		notify("Resource deleted", query)

		log.Infof("Done in %s", time.Now().Sub(startTime))

		return nil
	},
}

func deleteResource(ctx context.Context, urlStr string) error {
	body := &dflimg.IdentifyResource{
		Query: urlStr,
	}

	return makeClient().DeleteResource(ctx, body)
}

func handleQueryInput(args []string) (string, error) {
	if len(args) == 1 {
		return strings.TrimPrefix(args[0], rootURL()), nil
	}

	query, err := queryPrompt.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(query, rootURL()), nil
}
