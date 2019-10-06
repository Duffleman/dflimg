package cmd

import (
	"bytes"
	"dflimg"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dflimg/cmd/dflimg/http"

	"github.com/atotto/clipboard"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var UploadCmd = &cobra.Command{
	Use:     "upload",
	Aliases: []string{"u"},
	Short:   "Upload a file",
	Long:    "Upload a file from your local machine to a dflimg server",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		startTime := time.Now()

		localFile := args[0]
		shortcuts := cmd.Flag("shortcuts")
		nsfw := cmd.Flag("nsfw")

		rootURL := viper.Get("ROOT_URL").(string)
		authToken := viper.Get("AUTH_TOKEN").(string)

		body, err := sendFile(rootURL, authToken, localFile, shortcuts, nsfw)
		if err != nil {
			return err
		}

		err = clipboard.WriteAll(body.URL)
		if err != nil {
			log.Warn("Could not copy to clipboard. Please copy the URL manually")
		}
		notify("Image uploaded", body.URL)

		duration := time.Now().Sub(startTime)

		log.Infof("Done in %s: %s\n", duration, body.URL)

		return nil
	},
}

// SendFile uploads the file to the server
func sendFile(rootURL, authToken, filename string, shortcuts, nsfw *pflag.Flag) (*dflimg.CreateResourceResponse, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

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

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, err
	}

	io.Copy(part, file)
	writer.Close()

	c := http.New(rootURL, authToken)

	res := &dflimg.CreateResourceResponse{}
	err = c.Request("POST", "upload_file", body, writer.FormDataContentType(), res)

	return res, err
}
