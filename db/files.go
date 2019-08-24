package db

import (
	"context"
	"database/sql"

	"dflimg"
	"dflimg/dflerr"

	"github.com/lib/pq"
)

// NewFile inserts a new file into the database
func (db *DB) NewFile(ctx context.Context, id, s3, mimetype, owner string, shortcuts []string) (*dflimg.Resource, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Insert("resources").
		Columns("id, type, owner, link, mime_type, shortcuts").
		Values(id, "file", owner, s3, mimetype, pq.Array(shortcuts)).
		Suffix("RETURNING id, type, serial, owner, link, nsfw, mime_type, shortcuts, created_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	row := db.pg.QueryRowContext(ctx, query, values...)

	res := &dflimg.Resource{}

	err = row.Scan(
		&res.ID,
		&res.Type,
		&res.Serial,
		&res.Owner,
		&res.Link,
		&res.NSFW,
		&res.MimeType,
		pq.Array(&res.Shortcuts),
		&res.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, dflerr.New(dflerr.NotFound, nil)
		}
		return nil, err
	}

	return res, nil
}
