package repository

import (
	"context"

	"github.com/sinulingga23/dynamic-form/backend/model"
)

type IPFieldTypeRepository interface {
	FindOne(ctx context.Context, id string) (model.PFieldType, error)
	FindPFieldTypesByIds(ctx context.Context, ids []string) ([]model.PFieldType, error)
	Create(ctx context.Context, createPFieldType model.CreatePFieldType) error
	Creates(ctx context.Context, createPFieldTypes []model.CreatePFieldType) error
}
