package db

import (
	"context"

	"dflimg"

	"github.com/kr/pretty"
)

func (db *DB) ListLabels(ctx context.Context) ([]*dflimg.Label, error) {
	b := NewQueryBuilder()

	query, _, err := b.Select("id, name").From("labels").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.pg.QueryContext(ctx, query)
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

		pretty.Println(l)

		labels = append(labels, l)
	}

	return labels, nil
}
