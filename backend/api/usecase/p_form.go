package usecase

import (
	"context"

	"github.com/sinulingga23/dynamic-form/backend/payload"
)

type IPFormUsecase interface {
	AddPForm(ctx context.Context, pFormRequest payload.PFormRequest) payload.Response
	GetPFormsByPPartnerId(ctx context.Context, partnerId string) payload.Response
	GetPFormById(ctx context.Context, id string) payload.Response
}
