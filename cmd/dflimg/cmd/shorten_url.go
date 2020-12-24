package cmd

import (
	"context"
	"fmt"
	"time"

	"dflimg"

	"github.com/atotto/clipboard"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ShortenURLCmd = &cobra.Command{
	Use:     "shorten {url}",
	Aliases: []string{"s"},
	Short:   "Shorten a URL",
	Long:    "Shorten a URL",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		startTime := time.Now()

		urlStr := args[0]

		body, err := makeClient().ShortenURL(ctx, &dflimg.CreateURLRequest{
			URL: urlStr,
		})
		if err != nil {
			return err
		}

		err = clipboard.WriteAll(body.URL)
		if err != nil {
			fmt.Println("Could not copy to clipboard. Please copy the URL manually")
		}
		notify("URL Shortened", body.URL)

		log.Infof("Done in %s: %s", time.Now().Sub(startTime), body.URL)

		return nil
	},
}
