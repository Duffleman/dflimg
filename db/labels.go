package db

import (
	"context"

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

func (db *DB) queryLabels(ctx context.Context, query string, values []interface{}) ([]*dflimg.Label, error) {
	rows, err := db.pg.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	labels := []*dflimg.Label{}

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

	_, err = db.pg.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}
