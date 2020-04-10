package cmd

import (
	"time"

	"dflimg"
	"dflimg/cmd/dflimg/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DeleteResourceCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d"},
	Short:   "Delete a resource",
	Long:    "Delete a resource",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		startTime := time.Now()

		urlStr := args[0]

		rootURL := viper.Get("ROOT_URL").(string)
		authToken := viper.Get("AUTH_TOKEN").(string)

		err := deleteResource(rootURL, authToken, urlStr)
		if err != nil {
			return err
		}
		notify("Resource deleted", urlStr)

		duration := time.Now().Sub(startTime)

		log.Infof("Done in %s\n", duration)

		return nil
	},
}

func deleteResource(rootURL, authToken, urlStr string) error {
	body := &dflimg.IdentifyResource{
		Query: urlStr,
	}

	c := http.New(rootURL, authToken)

	return c.JSONRequest("POST", "delete_resource", body, nil)
}
