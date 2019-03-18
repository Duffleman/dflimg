package db

import (
	"time"

	"github.com/go-pg/pg"
)

// DB is a wrapper around the PG wrapper for easy function calls
type DB struct {
	pg *pg.DB
}

// New creates a new instance of a PG connection
func New(pg *pg.DB) *DB {
	return &DB{
		pg: pg,
	}
}

// File is a file DB entry
type File struct {
	ID        string    `sql:"id"`
	Serial    int       `sql:"serial"`
	Owner     string    `sql:"owner"`
	S3        string    `sql:"s3"`
	Type      string    `sql:"type"`
	CreatedAt time.Time `sql:"created_at"`
	Labels    []string  `sql:"labels,array"`
}

// NewFile inserts a new file into the database
func (db *DB) NewFile(id, s3, owner, t string, labels []string) error {
	inFile := &File{
		ID:     id,
		S3:     s3,
		Owner:  owner,
		Type:   t,
		Labels: labels,
	}

	return db.pg.Insert(inFile)
}

// FindFile retrieves a file from the database by it's ID
func (db *DB) FindFile(id string) (*File, error) {
	f := &File{
		ID: id,
	}

	err := db.pg.Select(f)

	return f, err
}

// FindFileBySerial retireves a file by it's serial number
func (db *DB) FindFileBySerial(serial int) (*File, error) {
	file := &File{}
	err := db.pg.Model(file).
		Where("serial = ?", serial).
		Limit(1).
		Select()

	return file, err
}

// FindFileByLabel retireves a file by a label
func (db *DB) FindFileByLabel(label string) (*File, error) {
	file := &File{}
	err := db.pg.Model(file).
		Where("? = ANY(labels)", label).
		Limit(1).
		Select()

	return file, err
}
