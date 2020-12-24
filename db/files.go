package db

import (
	"context"
	"fmt"
	"strings"

	"dflimg"
)

// NewFile inserts a new file into the database
func (db *DB) NewFile(ctx context.Context, id, s3, owner string, name *string, mimetype string) (*dflimg.Resource, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Insert("resources").
		Columns("id, type, owner, name, link, mime_type").
		Values(id, "file", owner, name, s3, mimetype).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(resourceColumns, ","))).
		ToSql()
	if err != nil {
		return nil, err
	}

	return db.queryOne(ctx, query, values)
}

// NewPendingFile inserts a new pending file into the database
func (db *DB) NewPendingFile(ctx context.Context, id, s3, owner string, name *string, mimetype string) (*dflimg.Resource, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Insert("resources").
		Columns("id, type, owner, name, link, mime_type").
		Values(id, "file", owner, name, s3, mimetype).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(resourceColumns, ","))).
		ToSql()
	if err != nil {
		return nil, err
	}

	return db.queryOne(ctx, query, values)
}
