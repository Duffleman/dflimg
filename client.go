package dflimg

import (
	"context"
	"net/http"
	"time"

	"dflimg/lib/crpc"
	"dflimg/lib/jsonclient"
)

type client struct {
	*crpc.Client
}

func NewClient(baseURL, key string) Service {
	httpClient := &http.Client{
		Transport: jsonclient.NewAuthenticatedRoundTripper(nil, key),
		Timeout:   5 * time.Second,
	}

	return &client{
		crpc.NewClient(baseURL+"/", httpClient),
	}
}

func (c *client) AddShortcut(ctx context.Context, req *ChangeShortcutRequest) error {
	return c.Do(ctx, "add_shortcut", req, nil)
}

func (c *client) CreatedSignedURL(ctx context.Context, req *CreateSignedURLRequest) (res *CreateSignedURLResponse, err error) {
	return res, c.Do(ctx, "create_signed_url", req, &res)
}

func (c *client) DeleteResource(ctx context.Context, req *IdentifyResource) error {
	return c.Do(ctx, "delete_resource", req, nil)
}

func (c *client) ListResources(ctx context.Context, req *ListResourcesRequest) (res []*Resource, err error) {
	return res, c.Do(ctx, "list_resources", req, &res)
}

func (c *client) RemoveShortcut(ctx context.Context, req *ChangeShortcutRequest) error {
	return c.Do(ctx, "remove_shortcut", req, nil)
}

func (c *client) SetNSFW(ctx context.Context, req *SetNSFWRequest) error {
	return c.Do(ctx, "set_nsfw", req, nil)
}

func (c *client) ShortenURL(ctx context.Context, req *CreateURLRequest) (res *CreateResourceResponse, err error) {
	return res, c.Do(ctx, "shorten_url", req, &res)
}

func (c *client) ViewDetails(ctx context.Context, req *IdentifyResource) (res *Resource, err error) {
	return res, c.Do(ctx, "view_details", req, &res)
}
