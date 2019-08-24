package app

import (
	"context"
	"fmt"

	"dflimg"
)

func (a *App) CreateResource(ctx context.Context, resourceID string) (*dflimg.ResponseCreatedResponse, error) {
	res, err := a.db.FindResource(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	rootURL := dflimg.GetEnv("root_url")
	hash := a.makeHash(res.Serial)
	fullURL := fmt.Sprintf("%s/%s", rootURL, hash)

	return &dflimg.ResponseCreatedResponse{
		ResourceID: res.ID,
		Type:       res.Type,
		Hash:       hash,
		URL:        fullURL,
	}, err
}
