package repository

import (
	"database/sql"

	"github.com/sinulingga23/dynamic-form/backend/api/repository"
)

type pFormFieldRepositoryImpl struct {
	db *sql.DB
}

func NewPFormFieldRepositoryImpl(db *sql.DB) repository.IPFormFieldRepository {
	return &pFormFieldRepositoryImpl{db: db}
}
