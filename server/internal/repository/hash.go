package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"hashServer/internal/generated/models"
	"strconv"
	"strings"
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

func (h *Hash) GetByIds(ids []string) (models.ArrayOfHash, error) {
	var hashes models.ArrayOfHash

	for _, v := range ids {
		if _, err := strconv.Atoi(v); err != nil {
			return hashes, errors.New("id must be a number")
		}
	}

	query := fmt.Sprintf("SELECT id, hash FROM %s WHERE id IN (%s)", hashesTable, strings.Join(ids, ","))

	err := h.db.Select(&hashes, query)

	return hashes, err
}
