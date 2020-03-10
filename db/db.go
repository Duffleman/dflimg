package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

// DB is a wrapper around the PG wrapper for easy function calls
type DB struct {
	pg *pgx.Conn
}

// New creates a new instance of a PG connection
func New(pg *pgx.Conn) *DB {
	return &DB{
		pg: pg,
	}
}

// NewQueryBuilder returns a new query builder
func NewQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
