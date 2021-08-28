package repository

import "github.com/jmoiron/sqlx"

type Hash struct {
	db *sqlx.DB
}

func NewHashPostgres(db *sqlx.DB) *Hash {
	return &Hash{db: db}
}
