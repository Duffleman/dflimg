package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DB is a wrapper around the PG wrapper for easy function calls
type DB struct {
	pg *pgxpool.Pool
}

// New creates a new instance of a PG connection
func New(pg *pgxpool.Pool) *DB {
	return &DB{
		pg: pg,
	}
}

// NewQueryBuilder returns a new query builder
func NewQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
