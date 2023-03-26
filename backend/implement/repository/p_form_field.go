package repository

import (
	"context"
	"database/sql"

	"github.com/sinulingga23/dynamic-form/backend/api/repository"
	"github.com/sinulingga23/dynamic-form/backend/model"
)

type pFormFieldRepositoryImpl struct {
	db *sql.DB
}

func NewPFormFieldRepositoryImpl(db *sql.DB) repository.IPFormFieldRepository {
	return &pFormFieldRepositoryImpl{db: db}
}

func (p *pFormFieldRepositoryImpl) FindOne(ctx context.Context, id string) (model.PFormField, error) {
	return model.PFormField{}, nil
}

func (p *pFormFieldRepositoryImpl) FindPFormFieldsByIds(ctx context.Context, ids []string) ([]model.PFormField, error) {
	return []model.PFormField{}, nil
}

func (p *pFormFieldRepositoryImpl) Create(ctx context.Context, createPFormField model.PFormField) error {
	return nil
}

func (p *pFormFieldRepositoryImpl) Creates(ctx context.Context, createPFormFields []model.PFormField) error {
	return nil
}
