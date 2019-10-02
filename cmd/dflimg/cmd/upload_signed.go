package cmd

import (
	"bytes"
	"dflimg"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	dhttp "dflimg/cmd/dflimg/http"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var UploadSignedCmd = &cobra.Command{
	Use:   "signed-upload",
	Short: "Upload a file to a signed URL",
	Long:  "Upload a file from your local machine to AWS",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		startTime := time.Now()

		localFile := args[0]

		rootURL := viper.Get("ROOT_URL").(string)
		authToken := viper.Get("AUTH_TOKEN").(string)

		body, err := sendFileAWS(rootURL, authToken, localFile)
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

// SendFileAWS uploads the file to AWS
func sendFileAWS(rootURL, authToken, filename string) (*dflimg.CreateSignedURLResponse, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(file)

	reqBody := &dflimg.CreateSignedURLRequest{
		ContentType: contentType,
	}

	c := dhttp.New(rootURL, authToken)

	res := &dflimg.CreateSignedURLResponse{}
	err = c.JSONRequest("POST", "created_signed_url", reqBody, res)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", res.S3Link, bytes.NewReader(file))
	if err != nil {
		return nil, err
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
