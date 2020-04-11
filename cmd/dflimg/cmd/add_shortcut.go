package cmd

import (
	"fmt"
	"time"

	"dflimg"
	"dflimg/cmd/dflimg/http"

	"github.com/atotto/clipboard"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var AddShortcutCmd = &cobra.Command{
	Use:     "add-shortcut",
	Aliases: []string{"asc"},
	Short:   "Add a shortcut",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		startTime := time.Now()

		query := args[0]
		shortcut := args[1]

		rootURL := viper.Get("ROOT_URL").(string)
		authToken := viper.Get("AUTH_TOKEN").(string)

		err := addShortcut(rootURL, authToken, query, shortcut)
		if err != nil {
			return err
		}

		err = clipboard.WriteAll(fmt.Sprintf("%s/:%s", rootURL, shortcut))
		if err != nil {
			log.Warn("Could not copy to clipboard.")
		}

		duration := time.Now().Sub(startTime)

		log.Infof("Done in %s", duration)

		return nil
	},
}

func addShortcut(rootURL, authToken, query, shortcut string) error {
	body := &dflimg.ChangeShortcutRequest{
		IdentifyResource: dflimg.IdentifyResource{
			Query: query,
		},
		Shortcut: shortcut,
	}

	c := http.New(rootURL, authToken)

	return c.JSONRequest("POST", "add_shortcut", body, nil)
}
