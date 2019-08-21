package db

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

// DB is a wrapper around the PG wrapper for easy function calls
type DB struct {
	pg *sql.DB
}

// New creates a new instance of a PG connection
func New(pg *sql.DB) *DB {
	return &DB{
		pg: pg,
	}
}

// NewQueryBuilder returns a new query builder
func NewQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
