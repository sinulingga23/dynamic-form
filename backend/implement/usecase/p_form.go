package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/sinulingga23/dynamic-form/backend/api/repository"
	"github.com/sinulingga23/dynamic-form/backend/api/usecase"
	"github.com/sinulingga23/dynamic-form/backend/payload"
)

type pFormUsecaseImpl struct {
	pFormRepository repository.IPFormRepository
}

func NewPFormUsecasseImpl(
	pFormRepository repository.IPFormRepository,
) usecase.IPFormUsecase {
	return &pFormUsecaseImpl{pFormRepository: pFormRepository}
}

func (p *pFormUsecaseImpl) GetFormsByPartnerId(ctx context.Context, partnerId string) ([]*payload.FormPartnerResponse, error) {
	pFormsPartner, errFindFormsByPartnerId := p.pFormRepository.FindPFormsByPartnerId(ctx, partnerId)
	if errFindFormsByPartnerId != nil {
		return []*payload.FormPartnerResponse{}, errFindFormsByPartnerId
	}

	pFormsPartnerResponse := make([]*payload.FormPartnerResponse, 0)

	lenPFormsPartner := len(pFormsPartner)
	for i := 0; i < lenPFormsPartner; i++ {
		pFormPartner := pFormsPartner[i]

		updateAt := time.Time{}
		if pFormPartner.UpdatedAt.Valid {
			updateAt = pFormPartner.UpdatedAt.Time
		}

		pFormsPartnerResponse = append(pFormsPartnerResponse, &payload.FormPartnerResponse{
			Id:           pFormPartner.Id,
			Name:         pFormPartner.Name,
			PPartnerId:   pFormPartner.PPartnerId,
			PPartnerName: pFormPartner.PPartnerName,
			CreatedAt:    pFormPartner.CreatedAt,
			UpdatedAt:    updateAt,
		})
	}

	if len(pFormsPartnerResponse) == 0 {
		return []*payload.FormPartnerResponse{}, errors.New("data not found.")
	}

	return pFormsPartnerResponse, nil
}
