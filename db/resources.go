package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"dflimg"
	"dflimg/dflerr"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
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
		Where("shortcuts @> $1", pq.Array(shortcuts)).
		Limit(1).
		ToSql()

	row := db.pg.QueryRowContext(ctx, query, values...)

	var id string
	err = row.Scan(&id)

	if err == sql.ErrNoRows {
		return nil
	}

	return errors.New("shortcut conflict")
}

// FindResource retrieves a resource from the database by it's ID
func (db *DB) FindResource(ctx context.Context, id string) (*dflimg.Resource, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Select("id, type, serial, owner, link, nsfw, mime_type, shortcuts, created_at").
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
		Select("id, type, serial, owner, link, nsfw, mime_type, shortcuts, created_at").
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
		Select("r.id, r.type, r.serial, r.owner, r.link, r.nsfw, r.mime_type, r.shortcuts, r.created_at").
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
	row := db.pg.QueryRowContext(ctx, query, values...)

	res := &dflimg.Resource{}

	err := row.Scan(
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

// GetLabelsBySerial returns labels associated with a resource
func (db *DB) GetLabelsBySerial(ctx context.Context, serial int) ([]*dflimg.Label, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Select("l.id, l.name").
		From("resources r").
		Join("labels_resources lr ON lr.resource_id = r.id").
		Join("labels l ON l.id = lr.label_id").
		Where(sq.Eq{
			"r.serial": serial,
		}).
		ToSql()
	if err != nil {
		return nil, err
	}

	return db.queryLabels(ctx, query, values)
}

// GetLabelsByShortcut returns labels associated with a resource
func (db *DB) GetLabelsByShortcut(ctx context.Context, shortcut string) ([]*dflimg.Label, error) {
	b := NewQueryBuilder()

	s := fmt.Sprintf("{%s}", shortcut[1:])

	query, values, err := b.
		Select("l.id, l.name").
		From("resources r").
		Join("labels_resources lr ON lr.resource_id = r.id").
		Join("labels l ON l.id = lr.label_id").
		Where("r.shortcuts @> $1::text[]", s).
		ToSql()
	if err != nil {
		return nil, err
	}

	return db.queryLabels(ctx, query, values)
}
