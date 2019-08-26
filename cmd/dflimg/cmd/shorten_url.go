package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"dflimg"
	"dflimg/cmd/dflimg/http"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var ShortenURLCmd = &cobra.Command{
	Use:     "shorten",
	Aliases: []string{"s"},
	Short:   "Shorten a URL",
	Long:    "Shorten a URL",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		startTime := time.Now()

		urlStr := args[0]
		shortcuts := cmd.Flag("shortcuts")
		nsfw := cmd.Flag("nsfw")

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

		duration := time.Now().Sub(startTime)

		fmt.Printf("Done in %s: %s\n", duration, body.URL)

		return nil
	},
}

func shortenURL(rootURL, authToken, urlStr string, shortcuts, nsfw *pflag.Flag) (*dflimg.CreateResourceResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if shortcuts != nil {
		shortcutsStr := shortcuts.Value.String()
		part, err := writer.CreateFormField("shortcuts")
		if err != nil {
			return nil, err
		}

		io.Copy(part, strings.NewReader(shortcutsStr))
	}

	if nsfw != nil {
		nsfwStr := nsfw.Value.String()
		part, err := writer.CreateFormField("nsfw")
		if err != nil {
			return nil, err
		}

		io.Copy(part, strings.NewReader(nsfwStr))
	}

	part, err := writer.CreateFormField("url")
	if err != nil {
		return nil, err
	}

	io.Copy(part, strings.NewReader(urlStr))
	writer.Close()

	c := http.New(rootURL, authToken)

	res := &dflimg.CreateResourceResponse{}
	err = c.Request("POST", "shorten_url", body, writer.FormDataContentType(), res)

	return res, err
}
