package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Hash struct {
	db *sqlx.DB
}

func NewHashPostgres(db *sqlx.DB) *Hash {
	return &Hash{db: db}
}

func (h *Hash) Create(hash string) (int, error) {
	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (hash) values ($1) RETURNING id", hashesTable)

	row := h.db.QueryRow(createItemQuery, hash)
	err := row.Scan(&itemId)
	if err != nil {
		return 0, err
	}

	return itemId, nil
}
