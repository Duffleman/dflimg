package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"dflimg"
	"dflimg/lib/cher"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

type ArrayOperation string

const (
	ArrayAdd    ArrayOperation = "array_append"
	ArrayRemove ArrayOperation = "array_remove"
)

// resourceColumns is the set of columns to populate into the struct
var resourceColumns = []string{"id", "type", "serial", "hash", "name", "owner", "link", "nsfw", "mime_type", "shortcuts", "created_at", "deleted_at"}

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

	s := fmt.Sprintf("{%s}", shortcut)

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

// FindResourceByName retrieves a resource from the database by an exact name match
func (db *DB) FindResourceByName(ctx context.Context, name string) (*dflimg.Resource, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Select(strings.Join(resourceColumns, ",")).
		From("resources").
		Where(sq.Eq{
			"name": name,
		}).
		OrderBy("serial DESC").
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
		&res.Name,
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
			return nil, cher.New(cher.NotFound, nil)
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

// SaveHash saves the hash of a resource into the DB
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

func (db *DB) ListResources(ctx context.Context, req *dflimg.ListResourcesRequest) ([]*dflimg.Resource, error) {
	b := NewQueryBuilder()

	builder := b.
		Select(strings.Join(resourceColumns, ",")).
		From("resources")

	if !req.IncludeDeleted {
		builder = builder.Where(sq.Eq{"deleted_at": nil})
	}

	if req.Username != nil {
		builder = builder.Where(sq.Eq{"owner": *req.Username})
	}

	if req.FilterMime != nil {
		builder = builder.Where(sq.Like{"mime_type": fmt.Sprintf("%s%%", *req.FilterMime)})
	}

	builder = builder.OrderBy("created_at DESC")

	if req.Limit != nil {
		builder = builder.Limit(*req.Limit)
	}

	query, values, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	fmt.Println(query)

	conn, err := db.pg.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	resources := []*dflimg.Resource{}

	for rows.Next() {
		o := &dflimg.Resource{}

		err := rows.Scan(&o.ID, &o.Type, &o.Serial, &o.Hash, &o.Name, &o.Owner, &o.Link, &o.NSFW, &o.MimeType, &o.Shortcuts, &o.CreatedAt, &o.DeletedAt)
		if err != nil {
			return nil, err
		}

		resources = append(resources, o)
	}

	return resources, nil
}

// ListResourcesWithoutHash lists all resources where the hash is not saved
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

func (db *DB) ChangeShortcut(ctx context.Context, operation ArrayOperation, resourceID, shortcut string) error {
	b := NewQueryBuilder()

	query, values, err := b.
		Update("resources").
		Set("shortcuts", sq.Expr(fmt.Sprintf("%s(shortcuts, ?)", operation), shortcut)).
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
