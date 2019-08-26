package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"dflimg/dflerr"

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

		_, err := tagResource(rootURL, authToken, urlStr, tagsStr)
		if err != nil {
			return err
		}

		duration := time.Now().Sub(startTime)

		fmt.Printf("Done in %s\n", duration)

		return nil
	},
}

func tagResource(rootURL, authToken, urlStr, tagsStr string) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormField("url")
	if err != nil {
		return nil, err
	}

	io.Copy(part, strings.NewReader(urlStr))

	part, err = writer.CreateFormField("tags")
	if err != nil {
		return nil, err
	}

	io.Copy(part, strings.NewReader(tagsStr))

	writer.Close()

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/tag_resource", rootURL), body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", authToken)
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var dflE dflerr.E
		err := json.Unmarshal(content, &dflE)
		if err != nil {
			return nil, err
		}

		return nil, dflE
	}

	return content, nil
}
