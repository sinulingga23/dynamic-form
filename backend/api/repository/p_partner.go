package repository

import (
	"context"

	"github.com/sinulingga23/dynamic-form/backend/model"
)

type IPPartnerRepository interface {
	FindOne(ctx context.Context, id string) (model.PPartner, error)
	FIndPPartnersByIds(ctx context.Context, ids []string) ([]model.PPartner, error)
	Create(ctx context.Context, createPartner model.CreatePPartner) error
	IsExistsById(ctx context.Context, id string) (bool, error)
}
