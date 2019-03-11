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
	ID        string
	Serial    int
	Owner     string
	S3        string
	Type      string
	CreatedAt time.Time
}

func (db *DB) NewFile(id, owner, t string) error {
	inFile := &File{
		ID:    id,
		Owner: owner,
		Type:  t,
	}

	return db.pg.Insert(inFile)
}
