package cmd

import (
	"bytes"
	"dflimg"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	dhttp "dflimg/cmd/dflimg/http"

	"github.com/kr/pretty"
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

		pretty.Println(body)

		duration := time.Now().Sub(startTime)

		fmt.Printf("Done in %s\n", duration)

		return nil
	},
}

// SendFileAWS uploads the file to AWS
func sendFileAWS(rootURL, authToken, filename string) (*dflimg.CreateResourceResponse, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := make([]byte, 512)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(buffer)

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	contentLength := fi.Size()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormField("content-type")
	if err != nil {
		return nil, err
	}
	io.Copy(part, strings.NewReader(contentType))

	part, err = writer.CreateFormField("content-length")
	if err != nil {
		return nil, err
	}
	io.Copy(part, strings.NewReader(strconv.FormatInt(contentLength, 10)))

	writer.Close()

	c := dhttp.New(rootURL, authToken)

	res := &dflimg.CreateSignedURLResponse{}
	err = c.Request("POST", "created_signed_url", body, writer.FormDataContentType(), res)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", res.URL, file)
	if err != nil {
		return nil, err
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
