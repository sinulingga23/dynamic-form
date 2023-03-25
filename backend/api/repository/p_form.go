package repository

import (
	"context"

	"github.com/sinulingga23/dynamic-form/backend/model"
)

type IPFormRepository interface {
	FindOne(ctx context.Context, id string) (model.PForm, error)
	FindPFormsByIds(ctx context.Context, ids []string) ([]model.PForm, error)
	Create(ctx context.Context, createPForm model.CreatePForm) error
}
