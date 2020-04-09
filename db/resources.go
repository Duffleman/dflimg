package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"dflimg"
	"dflimg/dflerr"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

// resourceColumns is the set of columns to populate into the struct
var resourceColumns = []string{"id", "type", "serial", "hash", "owner", "link", "nsfw", "mime_type", "shortcuts", "created_at", "deleted_at"}

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

	conn, err := db.pg.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	var id string

	err = conn.QueryRow(ctx, query, values...).Scan(&id)

	if err == pgx.ErrNoRows {
		return nil
	}

	return errors.New("shortcut conflict")
}

// FindResourceByHash retrieves a resource from the database by it's hash
func (db *DB) FindResourceByHash(ctx context.Context, hash string) (*dflimg.Resource, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Select(strings.Join(resourceColumns, ",")).
		From("resources").
		Where(sq.Eq{
			"hash": hash,
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
		Select(strings.Join(resourceColumns, ", ")).
		From("resources").
		Where("shortcuts @> $1::text[]", s).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	return db.queryOne(ctx, query, values)
}

func (db *DB) queryOne(ctx context.Context, query string, values []interface{}) (*dflimg.Resource, error) {
	res := &dflimg.Resource{}

	conn, err := db.pg.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(ctx, query, values...).Scan(
		&res.ID,
		&res.Type,
		&res.Serial,
		&res.Hash,
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

	conn, err := db.pg.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, query, values...)
	return err
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

	conn, err := db.pg.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, query, values...)
	return err
}

func (db *DB) SaveHash(ctx context.Context, serial int, hash string) error {
	b := NewQueryBuilder()

	query, values, err := b.
		Update("resources").
		Set("hash", hash).
		Where(sq.Eq{"serial": serial}).
		ToSql()
	if err != nil {
		return err
	}

	conn, err := db.pg.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, query, values...)
	return err
}

func (db *DB) ListResourcesWithoutHash(ctx context.Context) ([]*dflimg.ShortFormResource, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Select("id, serial").
		From("resources").
		Where(sq.Eq{
			"hash": nil,
		}).
		ToSql()

	conn, err := db.pg.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	resources := []*dflimg.ShortFormResource{}

	for rows.Next() {
		o := &dflimg.ShortFormResource{}

		err := rows.Scan(&o.ID, &o.Serial)
		if err != nil {
			return nil, err
		}

		resources = append(resources, o)
	}

	return resources, nil
}
