package db

import (
	"context"

	"dflimg"

	"github.com/lib/pq"
)

// NewFile inserts a new file into the database
func (db *DB) NewFile(ctx context.Context, id, s3, mimetype, owner string, shortcuts []string, nsfw bool) (*dflimg.Resource, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Insert("resources").
		Columns("id, type, owner, link, mime_type, shortcuts, nsfw").
		Values(id, "file", owner, s3, mimetype, pq.Array(shortcuts), nsfw).
		Suffix("RETURNING id, type, serial, owner, link, nsfw, mime_type, shortcuts, created_at, deleted_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	return db.queryOne(ctx, query, values)
}

// NewPendingFile inserts a new pending file into the database
func (db *DB) NewPendingFile(ctx context.Context, id, s3, mimetype, owner string, shortcuts []string, nsfw bool) (*dflimg.Resource, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Insert("resources").
		Columns("id, type, owner, link, mime_type, shortcuts, nsfw, pending").
		Values(id, "file", owner, s3, mimetype, pq.Array(shortcuts), nsfw, true).
		Suffix("RETURNING id, type, serial, owner, link, nsfw, mime_type, shortcuts, created_at, deleted_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	return db.queryOne(ctx, query, values)
}
