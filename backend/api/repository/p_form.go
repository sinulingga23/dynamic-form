package repository

import (
	"context"

	"github.com/sinulingga23/dynamic-form/backend/model"
)

type IPFormRepository interface {
	FindPFormsByPartnerId(ctx context.Context, partnerId string) ([]*model.PFormPartner, error)
	FindOne(ctx context.Context, id string) (*model.PFormDetail, error)
}
