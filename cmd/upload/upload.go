package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dflimg"

	"github.com/atotto/clipboard"
)

const (
	// URL of the dflimg server
	URL = "https://dfl.mn"
)

func main() {
	labelStr := flag.String("labels", "", "CSV list of labels")

	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		fin(errors.New("expecting exactly 1 file input"))
	}

	file := args[0]

	authToken, err := getAuthorisationToken()
	if err != nil {
		fin(err)
	}

	startTime := time.Now()

	body, err := sendFile(authToken, file, *labelStr)
	if err != nil {
		fin(err)
	}

	f, err := parseResponse(body)
	if err != nil {
		fin(err)
	}

	duration := time.Now().Sub(startTime)

	clipboard.WriteAll(f.URL)

	fmt.Printf("Done in %s: %s\n", duration.String(), f.URL)
}

func getAuthorisationToken() (string, error) {
	v := os.Getenv("DFLIMG_AUTH_TOKEN")
	if v == "" {
		return "", errors.New("no auth token set in env variables")
	}

	return v, nil
}

func parseResponse(res []byte) (*dflimg.UploadFileResponse, error) {
	var file dflimg.UploadFileResponse

	err := json.Unmarshal(res, &file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

// SendFile uploads the file to the server
func sendFile(authToken, filename, labelStr string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if labelStr != "" {
		part, err := writer.CreateFormField("labels")
		if err != nil {
			return nil, err
		}

		io.Copy(part, strings.NewReader(labelStr))
	}

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, err
	}

	io.Copy(part, file)
	writer.Close()

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/upload", URL), body)
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

	return content, nil
}

func fin(err error) {
	fmt.Println(err)

	os.Exit(1)
}
