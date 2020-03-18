package cmd

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"dflimg"
	dhttp "dflimg/cmd/dflimg/http"

	"github.com/atotto/clipboard"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func setup() (rootURL, authToken string) {
	rootURL = viper.Get("ROOT_URL").(string)
	authToken = viper.Get("AUTH_TOKEN").(string)

	return
}

var UploadSignedCmd = &cobra.Command{
	Use:     "signed-upload",
	Aliases: []string{"u"},
	Short:   "Upload a file to a signed URL",
	Long:    "Upload a file from your local machine to AWS",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		startTime := time.Now()

		rootURL, authToken := setup()

		localFile := args[0]
		shortcuts := cmd.Flag("shortcuts")
		nsfw := cmd.Flag("nsfw")

		file, err := ioutil.ReadFile(localFile)
		if err != nil {
			return err
		}

		filePrepStart := time.Now()
		resource, err := prepareUpload(rootURL, authToken, file, shortcuts, nsfw)
		if err != nil {
			return err
		}
		log.Infof("File prepared: %s (%s)\n", resource.URL, time.Now().Sub(filePrepStart))

		err = clipboard.WriteAll(resource.URL)
		if err != nil {
			log.Warn("Could not copy to clipboard. Please copy the URL manually")
		}
		notify("Image prepared", resource.URL)

		err = sendFileAWS(resource.S3Link, file)
		if err != nil {
			return err
		}
		notify("Image uploaded", "")

		log.Infof("Done in %s\n", time.Now().Sub(startTime))

		return nil
	},
}

func prepareUpload(rootURL, authToken string, file []byte, shortcuts, nsfw *pflag.Flag) (*dflimg.CreateSignedURLResponse, error) {
	contentType := http.DetectContentType(file)

	reqBody := &dflimg.CreateSignedURLRequest{
		ContentType: contentType,
	}

	if shortcuts != nil {
		shortcutsStr := shortcuts.Value.String()
		if shortcutsStr != "" {
			reqBody.Shortcuts = strings.Split(shortcutsStr, ",")
		}
	}

	if nsfw != nil {
		nsfwStr := nsfw.Value.String()

		if nsfwStr == "true" {
			reqBody.NSFW = true
		}
	}

	c := dhttp.New(rootURL, authToken)

	res := &dflimg.CreateSignedURLResponse{}
	err := c.JSONRequest("POST", "create_signed_url", reqBody, &res)

	return res, err
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
