package usecase

import (
	"context"

	"github.com/sinulingga23/dynamic-form/backend/payload"
)

type IPFormUsecase interface {
	GetFormsByPartnerId(ctx context.Context, partnerId string) ([]*payload.FormPartnerResponse, error)
}
