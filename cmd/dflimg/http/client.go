package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"dflimg/dflerr"
)

type Client struct {
	RootURL   string
	authToken string
	http      *http.Client
}

func New(rootURL, authToken string) *Client {
	c := http.DefaultClient

	return &Client{
		RootURL:   rootURL,
		authToken: authToken,
		http:      c,
	}
}

func (c *Client) Request(method, url string, body *bytes.Buffer, formType string, response interface{}) error {
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.RootURL, url), body)
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", formType)
	request.Header.Add("Authorization", c.authToken)

	res, err := c.http.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if len(content) == 0 {
		return nil
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var dflE dflerr.E
		err := json.Unmarshal(content, &dflE)
		if err != nil {
			return err
		}

		return dflE
	}

	err = json.Unmarshal(content, response)
	if err != nil {
		return err
	}

	return nil
}
