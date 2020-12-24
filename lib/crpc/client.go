package crpc

import (
	"context"
	"fmt"
	"net/http"

	"dflimg/lib/jsonclient"
)

const userAgentTemplate = "crpc (+https://cuvva.com)"

// Client represents a crpc client. It builds on top of jsonclient, so error
// variables/structs and the authenticated round tripper live there.
type Client struct {
	*jsonclient.Client
}

// NewClient returns a client configured with a transport scheme, remote host
// and URL prefix supplied as a URL <scheme>://<host></prefix>
func NewClient(baseURL string, c *http.Client) *Client {
	jcc := jsonclient.NewClient(baseURL, c)
	jcc.UserAgent = fmt.Sprintf(userAgentTemplate)

	return &Client{jcc}
}

// WithUASuffix updates the current user agent with an additional string
func (c *Client) WithUASuffix(suffix string) *Client {
	c.UserAgent = fmt.Sprintf("%s %s", c.UserAgent, suffix)
	return c
}

// Do executes an RPC request against the configured server.
func (c *Client) Do(ctx context.Context, method string, src, dst interface{}) error {
	err := c.Client.Do(ctx, "POST", method, nil, src, dst)

	if err == nil {
		return nil
	}

	if err, ok := err.(*jsonclient.ClientTransportError); ok {
		return &ClientTransportError{method, err.ErrorString, err.Cause()}
	}

	return err
}

// ClientTransportError is returned when an error related to
// executing a client request occurs.
type ClientTransportError struct {
	Method, ErrorString string

	cause error
}

// Cause returns the causal error (if wrapped) or nil
func (cte *ClientTransportError) Cause() error {
	return cte.cause
}

func (cte *ClientTransportError) Error() string {
	if cte.cause != nil {
		return fmt.Sprintf("%s %s: %s", cte.Method, cte.ErrorString, cte.cause.Error())
	}

	return fmt.Sprintf("%s %s", cte.Method, cte.ErrorString)
}
