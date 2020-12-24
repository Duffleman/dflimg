package cmd

import (
	"context"
	"time"

	"dflimg"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var DeleteResourceCmd = &cobra.Command{
	Use:     "delete {query}",
	Aliases: []string{"d"},
	Short:   "Delete a resource",
	Long:    "Delete a resource",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		startTime := time.Now()

		urlStr := args[0]

		err := deleteResource(ctx, urlStr)
		if err != nil {
			return err
		}
		notify("Resource deleted", urlStr)

		duration := time.Now().Sub(startTime)

		log.Infof("Done in %s", duration)

		return nil
	},
}

func deleteResource(ctx context.Context, urlStr string) error {
	body := &dflimg.IdentifyResource{
		Query: urlStr,
	}

	return makeClient().DeleteResource(ctx, body)
}
