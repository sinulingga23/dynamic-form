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

type pFormFieldUsecase struct {
	db                   *sql.DB
	pFormFieldRepository repository.IPFormFieldRepository
}

func NewPFormFieldUsecase(
	db *sql.DB,
	pFormFieldRepository repository.IPFormFieldRepository,
) usecase.IPFormFieldUsecase {
	return &pFormFieldUsecase{db: db, pFormFieldRepository: pFormFieldRepository}
}

func (p *pFormFieldUsecase) AddFormField(ctx context.Context, pFormFieldRequest payload.PFormFieldRequest) payload.Response {

	response := payload.Response{
		StatusCode: http.StatusOK,
		Message:    "Success add new form.",
	}

	if pFormFieldRequest.PFormName == "" {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Key formName cant be empty."
		return response
	}

	if pFormFieldRequest.PPartnerId == "" {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Key pPartnerId cant be empty."
		return response
	}

	_, errParse := uuid.Parse(pFormFieldRequest.PPartnerId)
	if errParse != nil {
		response.StatusCode = http.StatusNotFound
		response.Message = "Partner not found."
		return response
	}

	if len(pFormFieldRequest.FormFields) == 0 {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Key formFields cant be empty."
		return response
	}

	lenFormFields := len(pFormFieldRequest.FormFields)
	for i := 0; i < lenFormFields; i++ {
		formField := pFormFieldRequest.FormFields[i]

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
	rowQueryPartner := tx.QueryRow(queryCheckPartner, pFormFieldRequest.PPartnerId)

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
		pFormId, pFormFieldRequest.PPartnerId, pFormFieldRequest.PFormName, time.Now())
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
			paramFieldTypeIds += fmt.Sprintf(`'%s',`, pFormFieldRequest.FormFields[i].PFieldTypeId)
		} else {
			paramFieldTypeIds += fmt.Sprintf(`'%s'`, pFormFieldRequest.FormFields[i].PFieldTypeId)
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
		formField := pFormFieldRequest.FormFields[i]
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

func (p *pFormFieldUsecase) GetFormFieldById(ctx context.Context, id string) payload.Response {
	return payload.Response{}
}
