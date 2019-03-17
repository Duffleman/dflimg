package db

import (
	"time"

	"github.com/go-pg/pg"
)

type DB struct {
	pg *pg.DB
}

func New(pg *pg.DB) *DB {
	return &DB{
		pg: pg,
	}
}

type File struct {
	ID        string    `sql:"id"`
	Serial    int       `sql:"serial"`
	Owner     string    `sql:"owner"`
	S3        string    `sql:"s3"`
	Type      string    `sql:"type"`
	CreatedAt time.Time `sql:"created_at"`
}

func (db *DB) NewFile(id, s3, owner, t string) error {
	inFile := &File{
		ID:    id,
		S3:    s3,
		Owner: owner,
		Type:  t,
	}

	return db.pg.Insert(inFile)
}

func (db *DB) FindFile(id string) (*File, error) {
	f := &File{
		ID: id,
	}

	err := db.pg.Select(f)

	return f, err
}

func (db *DB) FindFileBySerial(serial int) (*File, error) {
	file := &File{}
	err := db.pg.Model(file).
		Where("serial = ?", serial).
		Limit(1).
		Select()

	return file, err
}
