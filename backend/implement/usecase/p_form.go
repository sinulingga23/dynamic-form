package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func NewPFormUsecaseImpl(
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

func (p *pFormUsecaseImpl) AddPForm(ctx context.Context, pFormRequest payload.PFormRequest) payload.Response {

	response := payload.Response{
		StatusCode: http.StatusOK,
		Message:    "Success add new form.",
	}

	if pFormRequest.PFormName == "" {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Key formName cant be empty."
		return response
	}

	if pFormRequest.PPartnerId == "" {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Key pPartnerId cant be empty."
		return response
	}

	_, errParse := uuid.Parse(pFormRequest.PPartnerId)
	if errParse != nil {
		response.StatusCode = http.StatusNotFound
		response.Message = "Partner not found."
		return response
	}

	if len(pFormRequest.PFormFields) == 0 {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Key formFields cant be empty."
		return response
	}

	lenFormFields := len(pFormRequest.PFormFields)
	for i := 0; i < lenFormFields; i++ {
		formField := pFormRequest.PFormFields[i]

		if formField.PFieldTypeId == "" {
			response.StatusCode = http.StatusBadRequest
			response.Message = "Key pFieldTypeId cant be empty."
			return response
		}

		_, errParse := uuid.Parse(formField.PFieldTypeId)
		if errParse != nil {
			response.StatusCode = http.StatusBadRequest
			response.Message = "FieldType not found."
			return response
		}

		if formField.PFormFieldName == "" {
			response.StatusCode = http.StatusBadRequest
			response.Message = "Key pFieldName cant be empty."
			return response
		}

		if formField.PFormFieldElement == "" {
			response.StatusCode = http.StatusBadRequest
			response.Message = "Key pFormFieldElement cant be empty."
			return response
		}
	}

	// START: Transaction
	tx, errBegin := p.db.Begin()
	if errBegin != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Error when rollback database transaction."
			return response
		}
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Error when start database transaction."
		return response
	}

	queryCheckPartner := `
	select
		count(id)
	from partner.p_partner 
		where id = $1
	`
	rowQueryPartner := tx.QueryRow(queryCheckPartner, pFormRequest.PPartnerId)

	countPartner := 0
	if errScanQueryPartner := rowQueryPartner.Scan(&countPartner); errScanQueryPartner != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Error when rollback database transaction."
			return response
		}

		if errors.Is(errScanQueryPartner, sql.ErrNoRows) {
			log.Printf("Error partner not found.")
			response.StatusCode = http.StatusNotFound
			response.Message = "Partner not found."
			return response
		}

		log.Printf("errScanQueryPartner: %v", errScanQueryPartner)
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Error query data."
		return response
	}

	if countPartner != 1 {
		log.Printf("Error countPartner: %v", countPartner)
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollbacl: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Error when rollback database transaction."
			return response
		}

		if countPartner == 0 {
			response.StatusCode = http.StatusNotFound
			response.Message = "Partner not found."
			return response
		}

		response.StatusCode = http.StatusInternalServerError
		response.Message = "Error query data."
		return response
	}

	queryInsertPForm := `
	insert into partner.p_form
		(id, p_partner_id, name, created_at)
	values
		($1, $2, $3, $4)
	`
	pFormId := uuid.NewString()
	resultQueryInsertPForm, errExecQueryInsertForm := tx.Exec(queryInsertPForm,
		pFormId, pFormRequest.PPartnerId, pFormRequest.PFormName, time.Now())
	if errExecQueryInsertForm != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollbacl: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Error when rollback database transaction."
			return response
		}
		log.Printf("errExecQueryInsertForm: %v", errExecQueryInsertForm)
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Error query data."
		return response
	}

	rowsAffectedQueryInsertForm, errRowsAffectedQueryInsertForm := resultQueryInsertPForm.RowsAffected()
	if errRowsAffectedQueryInsertForm != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Error when rollback database transaction."
			return response
		}
		log.Printf("errRowsAffectedQueryInsertForm: %v", errRowsAffectedQueryInsertForm)
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Error query data."
		return response
	}
	if rowsAffectedQueryInsertForm != int64(1) {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Error when rollback database transaction."
			return response
		}
		log.Printf("rowsAffectedQueryInsertForm: %v, should 1", rowsAffectedQueryInsertForm)
		response.StatusCode = http.StatusBadRequest
		response.Message = "Failed to insert data."
		return response
	}

	paramFieldTypeIds := `(`
	for i := 0; i < lenFormFields; i++ {
		if i != lenFormFields-1 {
			paramFieldTypeIds += fmt.Sprintf(`'%s',`, pFormRequest.PFormFields[i].PFieldTypeId)
		} else {
			paramFieldTypeIds += fmt.Sprintf(`'%s'`, pFormRequest.PFormFields[i].PFieldTypeId)
		}
	}
	paramFieldTypeIds += `)`

	queryCheckFieldTypeByIds := `
	select
		count(id)
	from
		partner.p_field_type
	where
		id in
	`
	queryCheckFieldTypeByIds += paramFieldTypeIds
	rowQueryCheckFieldTypeByIds := tx.QueryRow(queryCheckFieldTypeByIds)

	countFieldTypeByIds := 0
	if errScanQueryCheckFieldTypeByIds := rowQueryCheckFieldTypeByIds.Scan(&countFieldTypeByIds); errScanQueryCheckFieldTypeByIds != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Error when rollback database transaction."
			return response
		}
		log.Printf("errScanQueryCheckFieldTypeByIds: %v", errScanQueryCheckFieldTypeByIds)
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Error query data."
		return response
	}

	if errQueryCheckFieldTypeByIds := rowQueryCheckFieldTypeByIds.Err(); errQueryCheckFieldTypeByIds != nil {
		log.Printf("errQueryCheckFieldTypeByIds: %v", errQueryCheckFieldTypeByIds)
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Error query data."
		return response
	}

	paramPFormFields := ``
	for i := 0; i < lenFormFields; i++ {
		formField := pFormRequest.PFormFields[i]
		if i != lenFormFields-1 {
			paramPFormFields += fmt.Sprintf(`('%s', '%s', '%s', '%s', '%s', '%v'),`,
				uuid.NewString(),
				pFormId,
				formField.PFieldTypeId,
				formField.PFormFieldName,
				formField.PFormFieldElement,
				time.Now().Format(time.RFC3339))
		} else {
			paramPFormFields += fmt.Sprintf(`('%s', '%s', '%s', '%s', '%s', '%v')`,
				uuid.NewString(),
				pFormId,
				formField.PFieldTypeId,
				formField.PFormFieldName,
				formField.PFormFieldElement,
				time.Now().Format(time.RFC3339))
		}
	}

	queryInsertFormFields := `
	insert into partner.p_form_field
		(id, p_form_id, p_field_type_id, name, element, created_at)
	values `
	queryInsertFormFields += paramPFormFields
	resultQueryInserFormFields, errQueryInserFormFields := tx.Exec(queryInsertFormFields)
	if errQueryInserFormFields != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Error when rollback database transaction."
			return response
		}
		log.Println("errQueryInserFormFields:", errQueryInserFormFields)
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Error query data."
		return response
	}

	rowsAffectedQueryInsertFormFields, errRowsAffectedQueryInserFormFields := resultQueryInserFormFields.RowsAffected()
	if errRowsAffectedQueryInserFormFields != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Error when rollback database transaction."
			return response
		}
		log.Printf("errRowsAffectedQueryInserFormFields: %v", errRowsAffectedQueryInserFormFields)
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Error query data."
		return response
	}

	if lenFormFields != int(rowsAffectedQueryInsertFormFields) {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Error when rollback database transaction."
			return response
		}
		log.Printf("Failed to insert data. lenFormFields: %v, rowsAffectedQueryInsertFormFields: %v", lenFormFields, rowsAffectedQueryInsertFormFields)
		response.StatusCode = http.StatusBadRequest
		response.Message = "Failed to insert data."
		return response
	}

	if errCommit := tx.Commit(); errCommit != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Error when rollback database transaction."
			return response
		}
		log.Printf("errCommit: %v", errCommit)
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Error commit database transaction."
		return response
	}
	// END: Transaction

	return response
}

func (p *pFormUsecaseImpl) GetPFormsByPPartnerId(ctx context.Context, partnerId string) payload.Response {
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
