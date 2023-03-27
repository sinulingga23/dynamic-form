package usecase

import (
	"context"

	"github.com/sinulingga23/dynamic-form/backend/payload"
)

type IPFormFieldUsecase interface {
	AddFormField(ctx context.Context, pFormFieldRequest payload.PFormFieldRequest) payload.Response
	GetFormFieldById(ctx context.Context, id string) payload.Response
}
