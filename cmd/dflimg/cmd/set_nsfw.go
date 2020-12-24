package cmd

import (
	"context"
	"time"

	"dflimg"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var SetNSFWCmd = &cobra.Command{
	Use:     "nsfw {query}",
	Aliases: []string{"n"},
	Short:   "Toggle the NSFW flag",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		startTime := time.Now()

		query := args[0]

		newState, err := toggleNSFW(ctx, query)
		if err != nil {
			return err
		}

		log.Infof("NSFW flag is now %s", newState)

		duration := time.Now().Sub(startTime)

		log.Infof("Done in %s", duration)

		return nil
	},
}

func toggleNSFW(ctx context.Context, query string) (string, error) {
	res, err := makeClient().ViewDetails(ctx, &dflimg.IdentifyResource{
		Query: query,
	})
	if err != nil {
		return "", err
	}

	newState := "ON"

	if res.NSFW {
		newState = "OFF"
	}

	return newState, makeClient().SetNSFW(ctx, &dflimg.SetNSFWRequest{
		IdentifyResource: dflimg.IdentifyResource{
			Query: query,
		},
		NSFW: !res.NSFW,
	})
}
