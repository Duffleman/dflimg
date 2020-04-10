package db

import (
	"context"
	"fmt"
	"strings"

	"dflimg"
)

// NewURL inserts a new URL to the database
func (db *DB) NewURL(ctx context.Context, id, owner, url string) (*dflimg.Resource, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Insert("resources").
		Columns("id, type, owner, link").
		Values(id, "url", owner, url).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(resourceColumns, ","))).
		ToSql()
	if err != nil {
		return nil, err
	}

	return db.queryOne(ctx, query, values)
}
