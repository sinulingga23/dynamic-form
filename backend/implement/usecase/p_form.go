package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sinulingga23/dynamic-form/backend/api/repository"
	"github.com/sinulingga23/dynamic-form/backend/api/usecase"
	"github.com/sinulingga23/dynamic-form/backend/payload"
)

type pFormUsecaseImpl struct {
	db                 *sql.DB
	pFormRepository    repository.IPFormRepository
	pPartnerRepository repository.IPPartnerRepository
}

func NewPFormUsecasseImpl(
	db *sql.DB,
	pFormRepository repository.IPFormRepository,
	pPartnerRepository repository.IPPartnerRepository,
) usecase.IPFormUsecase {
	return &pFormUsecaseImpl{
		db:                 db,
		pFormRepository:    pFormRepository,
		pPartnerRepository: pPartnerRepository,
	}
}

func (p *pFormUsecaseImpl) GetFormsByPartnerId(ctx context.Context, partnerId string) payload.Response {
	response := payload.Response{
		StatusCode: http.StatusOK,
		Message:    "Success to get the forms.",
	}

	if partnerId == "" {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Param partnerId cant be empty."
		return response
	}

	_, errParse := uuid.Parse(partnerId)
	if errParse != nil {
		log.Printf("errParse: %v, partnerId: %v", errParse, partnerId)
		response.StatusCode = http.StatusNotFound
		response.Message = "Partner not found."
		return response
	}

	isPartnerExists, errIsExists := p.pPartnerRepository.IsExistsById(ctx, partnerId)
	if errIsExists != nil {
		log.Printf("errIsExists: %v", errIsExists)
		if errors.Is(errIsExists, sql.ErrNoRows) {
			response.Data = http.StatusNotFound
			response.Message = "Partner not found"
			return response
		}
		response.Data = http.StatusInternalServerError
		response.Message = "Error query data"
		return response
	}

	if !isPartnerExists {
		log.Printf("isPartnerExists: %v", isPartnerExists)
		response.Data = http.StatusNotFound
		response.Message = "Partner not found"
		return response
	}

	pFormsPartner, errFindFormsByPartnerId := p.pFormRepository.FindPFormsByPartnerId(ctx, partnerId)
	if errFindFormsByPartnerId != nil {
		log.Printf("errFindFormsByPartnerId: %v, partnerId: %v", errFindFormsByPartnerId, partnerId)
		if errors.Is(errFindFormsByPartnerId, sql.ErrNoRows) {
			response.StatusCode = http.StatusNotFound
			response.Message = "Data not found"
			return response
		}

		response.StatusCode = http.StatusInternalServerError
		response.Message = "Error query data."
		return response
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
		log.Printf("pFormsPartnerResponse: %v", pFormsPartnerResponse)
		response.StatusCode = http.StatusNotFound
		response.Message = "Data not found."
		return response
	}

	response.Data = pFormsPartnerResponse
	return response
}
