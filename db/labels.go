package db

import (
	"context"
	"fmt"

	"dflimg"

	sq "github.com/Masterminds/squirrel"
)

// ListLabels returns a slice with every label
func (db *DB) ListLabels(ctx context.Context) ([]*dflimg.Label, error) {
	b := NewQueryBuilder()

	query, _, err := b.Select("id, name").From("labels").ToSql()
	if err != nil {
		return nil, err
	}

	return db.queryLabels(ctx, query, nil)
}

// GetLabelsByName returns a slice of labels with matching names
func (db *DB) GetLabelsByName(ctx context.Context, names []string) ([]*dflimg.Label, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Select("id, name").
		From("labels").
		Where(sq.Eq{
			"name": names,
		}).
		ToSql()
	if err != nil {
		return nil, err
	}

	return db.queryLabels(ctx, query, values)
}

// GetLabelsByHash returns labels associated with a resource
func (db *DB) GetLabelsByHash(ctx context.Context, hash string) ([]*dflimg.Label, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Select("l.id, l.name").
		From("resources r").
		Join("labels_resources lr ON lr.resource_id = r.id").
		Join("labels l ON l.id = lr.label_id").
		Where(sq.Eq{
			"r.hash": hash,
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

func (db *DB) queryLabels(ctx context.Context, query string, values []interface{}) ([]*dflimg.Label, error) {
	labels := []*dflimg.Label{}

	conn, err := db.pg.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		l := &dflimg.Label{}

		err := rows.Scan(&l.ID, &l.Name)
		if err != nil {
			return nil, err
		}

		labels = append(labels, l)
	}

	return labels, nil
}
