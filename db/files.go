package db

import (
	"context"

	"github.com/lib/pq"
)

// NewFile inserts a new file into the database
func (db *DB) NewFile(ctx context.Context, id, s3, owner, mimetype string, shortcuts []string) error {
	b := NewQueryBuilder()

	query, values, err := b.
		Insert("resources").
		Columns("id, type, owner, link, mime_type, shortcuts").
		Values(id, "file", owner, s3, mimetype, pq.Array(shortcuts)).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.pg.ExecContext(ctx, query, values...)

	return err
}
