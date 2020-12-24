package cmd

import (
	"context"
	"time"

	"dflimg"
	"dflimg/lib/cher"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var SetNSFWCmd = &cobra.Command{
	Use:     "nsfw {query}",
	Aliases: []string{"n"},
	Short:   "Toggle the NSFW flag",
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

		newState, err := toggleNSFW(ctx, query)
		if err != nil {
			return err
		}

		log.Infof("NSFW flag is now %s", newState)

		log.Infof("Done in %s", time.Now().Sub(startTime))

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
