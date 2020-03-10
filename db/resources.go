package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"dflimg"
	"dflimg/dflerr"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
)

// FindShortcutConflicts returns error if a shortcut is already taken
func (db *DB) FindShortcutConflicts(ctx context.Context, shortcuts []string) error {
	if len(shortcuts) == 0 {
		return nil
	}

	b := NewQueryBuilder()

	query, values, err := b.
		Select("id").
		From("resources").
		Where("shortcuts @> $1", shortcuts).
		Limit(1).
		ToSql()

	row := db.pg.QueryRow(ctx, query, values...)

	var id string
	err = row.Scan(&id)

	if err == pgx.ErrNoRows {
		return nil
	}

	return errors.New("shortcut conflict")
}

// FindResource retrieves a resource from the database by it's ID
func (db *DB) FindResource(ctx context.Context, id string) (*dflimg.Resource, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Select("id, type, serial, owner, link, nsfw, mime_type, shortcuts, created_at, deleted_at").
		From("resources").
		Where(sq.Eq{
			"id": id,
		}).
		ToSql()
	if err != nil {
		return nil, err
	}

	return db.queryOne(ctx, query, values)
}

// FindResourceBySerial retrieves a resource from the database by it's serial allocation
func (db *DB) FindResourceBySerial(ctx context.Context, serial int) (*dflimg.Resource, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Select("id, type, serial, owner, link, nsfw, mime_type, shortcuts, created_at, deleted_at").
		From("resources").
		Where(sq.Eq{
			"serial": serial,
		}).
		ToSql()
	if err != nil {
		return nil, err
	}

	return db.queryOne(ctx, query, values)
}

// FindResourceByShortcut retrieves a resource from the database by one of it's shortcuts
func (db *DB) FindResourceByShortcut(ctx context.Context, shortcut string) (*dflimg.Resource, error) {
	b := NewQueryBuilder()

	s := fmt.Sprintf("{%s}", shortcut[1:])

	query, values, err := b.
		Select("r.id, r.type, r.serial, r.owner, r.link, r.nsfw, r.mime_type, r.shortcuts, r.created_at, r.deleted_at").
		From("resources r").
		Where("r.shortcuts @> $1::text[]", s).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	return db.queryOne(ctx, query, values)
}

func (db *DB) queryOne(ctx context.Context, query string, values []interface{}) (*dflimg.Resource, error) {
	row := db.pg.QueryRow(ctx, query, values...)

	res := &dflimg.Resource{}

	err := row.Scan(
		&res.ID,
		&res.Type,
		&res.Serial,
		&res.Owner,
		&res.Link,
		&res.NSFW,
		&res.MimeType,
		&res.Shortcuts,
		&res.CreatedAt,
		&res.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, dflerr.New(dflerr.NotFound, nil)
		}
		return nil, err
	}

	return res, nil
}

// SetNSFW sets a resource NSFW bool
func (db *DB) SetNSFW(ctx context.Context, resourceID string, state bool) error {
	b := NewQueryBuilder()

	query, values, err := b.
		Update("resources").
		Set("nsfw", state).
		Where(sq.Eq{"id": resourceID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.pg.Exec(ctx, query, values...)
	return err
}

// TagResource tags a resource with a label, it is idempotant
func (db *DB) TagResource(ctx context.Context, resourceID string, tags []*dflimg.Label) error {
	b := NewQueryBuilder()

	builder := b.
		Insert("labels_resources").
		Columns("label_id, resource_id")

	for _, t := range tags {
		builder = builder.Values(t.ID, resourceID)
	}

	builder = builder.Suffix("ON CONFLICT (label_id, resource_id) DO NOTHING")

	query, values, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = db.pg.Exec(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

// DeleteResource soft-deletes a resource
func (db *DB) DeleteResource(ctx context.Context, resourceID string) error {
	b := NewQueryBuilder()

	query, values, err := b.
		Update("resources").
		Set("deleted_at", time.Now()).
		Where(sq.Eq{"id": resourceID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.pg.Exec(ctx, query, values...)
	return err
}
