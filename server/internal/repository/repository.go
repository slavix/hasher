package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Hash *Hash
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Hash: NewHashPostgres(db),
	}
}
