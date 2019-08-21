package db

import (
	"context"
	"database/sql"
	"dflimg"
	"errors"

	sq "github.com/Masterminds/squirrel"
)

// NewFile inserts a new file into the database
func (db *DB) NewFile(ctx context.Context, id, s3, owner, t string, shortcuts []string) error {
	b := NewQueryBuilder()

	query, values, err := b.
		Insert("files").
		Columns("id, s3, owner, type, shortcuts").
		Values(id, s3, owner, t, shortcuts).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.pg.ExecContext(ctx, query, values...)

	return err
}

// FindFile retrieves a file from the database by it's ID
func (db *DB) FindFile(ctx context.Context, id string) (file *dflimg.File, err error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Select("id, owner, s3, type, shortcuts, created_at").
		From("files").
		Where(sq.Eq{
			"id": id,
		}).
		ToSql()

	row := db.pg.QueryRowContext(ctx, query, values...)

	err = row.Scan(&file.ID, &file.Owner, &file.S3, &file.Type, &file.Shortcuts, &file.CreatedAt)

	return
}

// FindFileBySerial retireves a file by it's serial number
func (db *DB) FindFileBySerial(ctx context.Context, serial int) (file *dflimg.File, err error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Select("id, owner, s3, type, shortcuts, created_at").
		From("files").
		Where(sq.Eq{
			"serial": serial,
		}).
		ToSql()

	row := db.pg.QueryRowContext(ctx, query, values...)

	err = row.Scan(&file.ID, &file.Owner, &file.S3, &file.Type, &file.Shortcuts, &file.CreatedAt)

	return
}

// FindFileByShortcut retireves a file by one of it's shortcut
func (db *DB) FindFileByShortcut(ctx context.Context, shortcut string) (file *dflimg.File, err error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Select("id, owner, s3, type, shortcuts, created_at").
		From("files").
		Where("shortcuts @> $1", shortcut).
		ToSql()

	row := db.pg.QueryRowContext(ctx, query, values...)

	err = row.Scan(&file.ID, &file.Owner, &file.S3, &file.Type, &file.Shortcuts, &file.CreatedAt)

	return
}

// FindShortcutConflicts returns error if a shortcut is already taken
func (db *DB) FindShortcutConflicts(ctx context.Context, shortcuts []string) error {
	if len(shortcuts) == 0 {
		return nil
	}

	b := NewQueryBuilder()

	query, values, err := b.
		Select("id").
		From("files").
		Where("shortcuts @> ANY($1)", shortcuts).
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
