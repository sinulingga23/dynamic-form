package repository

import (
	"context"

	"github.com/sinulingga23/dynamic-form/backend/model"
)

type IPFormFieldRepository interface {
	FindOne(ctx context.Context, id string) (model.PFormField, error)
	FindPFormFieldsByIds(ctx context.Context, ids []string) ([]model.PFormField, error)
	Create(ctx context.Context, createPFormField model.PFormField) error
	Creates(ctx context.Context, createPFormFields []model.PFormField) error
}
