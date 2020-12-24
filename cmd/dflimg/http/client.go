package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"dflimg/lib/cher"
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

func (c *Client) JSONRequest(method, url string, body, response interface{}) error {
	jsonRaw, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.RootURL, url), bytes.NewBuffer(jsonRaw))
	if err != nil {
		return err
	}

	return c.doRequest(request, response)
}

func (c *Client) Request(method, url string, body *bytes.Buffer, formType string, response interface{}) error {
	request, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.RootURL, url), body)
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", formType)

	return c.doRequest(request, response)
}

func (c *Client) doRequest(request *http.Request, response interface{}) error {
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
		var c cher.E
		err := json.Unmarshal(content, &c)
		if err != nil {
			return err
		}

		return c
	}

	err = json.Unmarshal(content, response)
	if err != nil {
		return err
	}

	return nil
}
