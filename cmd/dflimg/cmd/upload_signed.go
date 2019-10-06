package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"dflimg"
	dhttp "dflimg/cmd/dflimg/http"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func setup() (rootURL, authToken string) {
	rootURL = viper.Get("ROOT_URL").(string)
	authToken = viper.Get("AUTH_TOKEN").(string)

	return
}

var UploadSignedCmd = &cobra.Command{
	Use:   "signed-upload",
	Short: "Upload a file to a signed URL",
	Long:  "Upload a file from your local machine to AWS",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		startTime := time.Now()

		rootURL, authToken := setup()

		localFile := args[0]

		file, err := ioutil.ReadFile(localFile)
		if err != nil {
			return err
		}

		filePrepStart := time.Now()
		resource, err := prepareUpload(rootURL, authToken, file)
		if err != nil {
			return err
		}
		fmt.Printf("File prepared: %s (%s)\n", resource.URL, time.Now().Sub(filePrepStart))

		err = clipboard.WriteAll(resource.URL)
		if err != nil {
			fmt.Printf("Could not copy to clipboard. Please copy the URL manually")
		}

		err = sendFileAWS(resource.S3Link, file)
		if err != nil {
			return err
		}

		fmt.Printf("Done in %s\n", time.Now().Sub(startTime))

		return nil
	},
}

func prepareUpload(rootURL, authToken string, file []byte) (res *dflimg.CreateSignedURLResponse, err error) {
	contentType := http.DetectContentType(file)

	reqBody := &dflimg.CreateSignedURLRequest{
		ContentType: contentType,
	}

	c := dhttp.New(rootURL, authToken)

	err = c.JSONRequest("POST", "create_signed_url", reqBody, &res)

	return
}

// SendFileAWS uploads the file to AWS
func sendFileAWS(signedURL string, file []byte) error {
	contentType := http.DetectContentType(file)

	req, err := http.NewRequest("PUT", signedURL, bytes.NewReader(file))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return err
}
