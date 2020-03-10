package db

import (
	"context"

	"dflimg"
)

// NewURL inserts a new URL to the database
func (db *DB) NewURL(ctx context.Context, id, url, owner string, shortcuts []string, nsfw bool) (*dflimg.Resource, error) {
	b := NewQueryBuilder()

	if len(shortcuts) < 1 {
		shortcuts = nil
	}

	query, values, err := b.
		Insert("resources").
		Columns("id, type, owner, link, shortcuts, nsfw").
		Values(id, "url", owner, url, shortcuts, nsfw).
		Suffix("RETURNING id, type, serial, owner, link, nsfw, mime_type, shortcuts, created_at, deleted_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	return db.queryOne(ctx, query, values)
}
