package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dflimg"
	"dflimg/lib/cher"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var UploadSignedCmd = &cobra.Command{
	Use:     "signed-upload {file}",
	Aliases: []string{"u"},
	Short:   "Upload a file to a signed URL",
	Long:    "Upload a file from your local machine to AWS",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 || len(args) == 0 {
			return nil
		}

		return cher.New("missing_arguments", nil)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		startTime := time.Now()

		localFile, err := handleLocalFileInput(args)
		if err != nil {
			return err
		}

		filePaths, err := scanDirectory(localFile)
		if err != nil {
			return err
		}

		if len(filePaths) == 0 {
			return cher.New("no_files", nil)
		}

		all := []string{}

		singleFile := len(filePaths) == 1

		for _, filename := range filePaths {
			log.Infof("Handling file: %s", filename)
			innerStart := time.Now()

			file, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}

			filePrepStart := time.Now()

			resource, err := prepareUpload(ctx, filename, file)
			if err != nil {
				return err
			}

			all = append(all, resource.Hash)

			log.Infof("File prepared: %s (%s)", resource.URL, time.Now().Sub(filePrepStart))

			if singleFile {
				writeClipboard(resource.URL)
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

		if !singleFile {
			jointURL := fmt.Sprintf("%s%s", rootURL(), strings.Join(all, ","))
			log.Infof("Download TAR at: %s", jointURL)
			writeClipboard(jointURL)
		}

		return nil
	},
}

func prepareUpload(ctx context.Context, filename string, file []byte) (*dflimg.CreateSignedURLResponse, error) {
	contentType := http.DetectContentType(file)

	var name *string

	if filename != "" {
		_, tmpName := filepath.Split(filename)
		name = &tmpName
	}

	return makeClient().CreatedSignedURL(ctx, &dflimg.CreateSignedURLRequest{
		ContentType: contentType,
		Name:        name,
	})
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

func handleLocalFileInput(args []string) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}

	file, err := filePrompt.Run()
	if err != nil {
		return "", err
	}

	return file, nil
}
