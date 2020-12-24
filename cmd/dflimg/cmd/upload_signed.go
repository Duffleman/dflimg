package cmd

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dflimg"
	dhttp "dflimg/cmd/dflimg/http"
	"dflimg/dflerr"

	"github.com/atotto/clipboard"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

		filePaths, err := scanDirectory(localFile)
		if err != nil {
			return err
		}

		if len(filePaths) == 0 {
			return dflerr.New("no_fies", nil)
		}

		singleFile := len(filePaths) == 1

		for _, filename := range filePaths {
			log.Infof("Handling file: %s", filename)
			innerStart := time.Now()

			file, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}

			filePrepStart := time.Now()
			resource, err := prepareUpload(rootURL, authToken, filename, file)
			if err != nil {
				return err
			}
			log.Infof("File prepared: %s (%s)", resource.URL, time.Now().Sub(filePrepStart))

			if singleFile {
				err = clipboard.WriteAll(resource.URL)
				if err != nil {
					log.Warn("Could not copy to clipboard. Please copy the URL manually")
				}
				notify("Image prepared", resource.URL)
			}

			err = sendFileAWS(resource.S3Link, file)
			if err != nil {
				return err
			}
			if singleFile {
				notify("Image uploaded", resource.URL)
			} else {
				log.Infof("File uploaded: %s", resource.URL)
			}

			log.Infof("File handled in %s", time.Now().Sub(innerStart))
		}

		log.Infof("Done in %s", time.Now().Sub(startTime))

		return nil
	},
}

func prepareUpload(rootURL, authToken string, filename string, file []byte) (*dflimg.CreateSignedURLResponse, error) {
	contentType := http.DetectContentType(file)

	var name *string

	if filename != "" {
		tmpName := strings.Split(filename, "/")
		name = &tmpName[len(tmpName)-1]
	}

	reqBody := &dflimg.CreateSignedURLRequest{
		ContentType: contentType,
		Name:        name,
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

func scanDirectory(rootFile string) (filePaths []string, err error) {
	root, err := os.Stat(rootFile)
	if err != nil {
		return nil, err
	}

	if !root.IsDir() {
		return []string{rootFile}, nil
	}

	err = filepath.Walk(rootFile, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		filePaths = append(filePaths, path)

		return nil
	})

	return
}
