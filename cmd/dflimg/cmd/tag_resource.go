package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"dflimg/cmd/dflimg/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var TagResourceCmd = &cobra.Command{
	Use:     "tag",
	Aliases: []string{"t"},
	Short:   "Tag a resource",
	Long:    "Tag a resource with a label",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		startTime := time.Now()

		urlStr := args[0]
		tagsStr := args[1]

		rootURL := viper.Get("ROOT_URL").(string)
		authToken := viper.Get("AUTH_TOKEN").(string)

		err := tagResource(rootURL, authToken, urlStr, tagsStr)
		if err != nil {
			return err
		}

		duration := time.Now().Sub(startTime)

		fmt.Printf("Done in %s\n", duration)

		return nil
	},
}

func tagResource(rootURL, authToken, urlStr, tagsStr string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormField("url")
	if err != nil {
		return err
	}

	io.Copy(part, strings.NewReader(urlStr))

	part, err = writer.CreateFormField("tags")
	if err != nil {
		return err
	}

	io.Copy(part, strings.NewReader(tagsStr))

	writer.Close()

	c := http.New(rootURL, authToken)

	return c.Request("POST", "tag_resource", body, writer.FormDataContentType(), nil)
}
