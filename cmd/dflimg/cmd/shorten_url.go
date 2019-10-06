package cmd

import (
	"dflimg"
	"dflimg/cmd/dflimg/http"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ShortenURLCmd = &cobra.Command{
	Use:     "shorten",
	Aliases: []string{"s"},
	Short:   "Shorten a URL",
	Long:    "Shorten a URL",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		urlStr := args[0]
		shortcutsPf := cmd.Flag("shortcuts")
		nsfwPf := cmd.Flag("nsfw")
		nsfw := false

		if strings.ToLower(nsfwPf.Value.String()) == "true" {
			nsfw = true
		}

		shortcutsStr := shortcutsPf.Value.String()
		shortcuts := strings.Split(shortcutsStr, ",")

		rootURL := viper.Get("ROOT_URL").(string)
		authToken := viper.Get("AUTH_TOKEN").(string)

		body, err := shortenURL(rootURL, authToken, urlStr, shortcuts, nsfw)
		if err != nil {
			return err
		}

		err = clipboard.WriteAll(body.URL)
		if err != nil {
			fmt.Println("Could not copy to clipboard. Please copy the URL manually")
		}
		notify("URL Shortened", body.URL)

		log.Infof("Done: %s\n", body.URL)

		return nil
	},
}

func shortenURL(rootURL, authToken, urlStr string, shortcuts []string, nsfw bool) (*dflimg.CreateResourceResponse, error) {
	body := &dflimg.CreateResourceRequest{
		Type:      "url",
		URL:       urlStr,
		Shortcuts: shortcuts,
		NSFW:      nsfw,
	}

	c := http.New(rootURL, authToken)

	res := &dflimg.CreateResourceResponse{}
	err := c.JSONRequest("POST", "shorten_url", body, res)

	return res, err
}
