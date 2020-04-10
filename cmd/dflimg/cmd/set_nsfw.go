package cmd

import (
	"time"

	"dflimg"
	"dflimg/cmd/dflimg/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var SetNSFWCmd = &cobra.Command{
	Use:     "nsfw",
	Aliases: []string{"n"},
	Short:   "Toggle the NSFW flag",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		startTime := time.Now()

		query := args[0]

		rootURL := viper.Get("ROOT_URL").(string)
		authToken := viper.Get("AUTH_TOKEN").(string)

		newState, err := toggleNSFW(rootURL, authToken, query)
		if err != nil {
			return err
		}

		log.Infof("NSFW flag is now %s", newState)

		duration := time.Now().Sub(startTime)

		log.Infof("Done in %s\n", duration)

		return nil
	},
}

func toggleNSFW(rootURL, authToken, query string) (string, error) {
	body := &dflimg.IdentifyResource{
		Query: query,
	}
	c := http.New(rootURL, authToken)
	res := &dflimg.Resource{}

	err := c.JSONRequest("POST", "view_details", body, &res)
	if err != nil {
		return "", err
	}

	swapTo := &dflimg.SetNSFWRequest{
		IdentifyResource: dflimg.IdentifyResource{
			Query: query,
		},
		NSFW: !res.NSFW,
	}

	newState := "ON"

	if res.NSFW == true {
		newState = "OFF"
	}

	return newState, c.JSONRequest("POST", "set_nsfw", swapTo, nil)
}
