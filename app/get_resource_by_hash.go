package app

import (
	"context"

	"dflimg"
	"dflimg/dflerr"
)

func (a *App) GetResourceByHash(ctx context.Context, hash string) (*dflimg.Resource, error) {
	serial, err := a.decodeHash(hash)
	if err != nil {
		return nil, err
	}

	return a.db.FindResourceBySerial(ctx, serial)
}

func (a *App) decodeHash(hash string) (int, error) {
	var set []int

	set, err := a.hasher.DecodeWithError(hash)
	if len(set) != 1 {
		return 0, dflerr.New("cannot decode hash", dflerr.M{"hash": hash}, dflerr.New("expecting single hashed item in body", nil))
	}

	return set[0], err
}
