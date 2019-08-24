package db

import (
	"context"

	"github.com/lib/pq"
)

// NewURL inserts a new URL to the database
func (db *DB) NewURL(ctx context.Context, id, url, owner string, shortcuts []string) error {
	b := NewQueryBuilder()

	query, values, err := b.
		Insert("resources").
		Columns("id, type, owner, link, shortcuts").
		Values(id, "url", owner, url, pq.Array(shortcuts)).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.pg.ExecContext(ctx, query, values...)

	return err
}
