package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
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

func (p *pFormFieldUsecase) AddFormField(ctx context.Context, pFormFieldRequest payload.PFormFieldRequest) error {
	if pFormFieldRequest.PFormName == "" {
		return errors.New("Key formName cant be empty.")
	}

	if pFormFieldRequest.PPartnerId == "" {
		return errors.New("Key pPartnerId cant be empty.")
	}

	_, errParse := uuid.Parse(pFormFieldRequest.PPartnerId)
	if errParse != nil {
		return errors.New("Key pPartnerId is not valid id.")
	}

	if len(pFormFieldRequest.FormFields) == 0 {
		return errors.New("Key formFields cant be empty.")
	}

	lenFormFields := len(pFormFieldRequest.FormFields)
	for i := 0; i < lenFormFields; i++ {
		formField := pFormFieldRequest.FormFields[i]

		if formField.PFieldTypeId == "" {
			return errors.New("Key pFieldTypeId cant be empty.")
		}

		_, errParse := uuid.Parse(formField.PFieldTypeId)
		if errParse != nil {
			return errors.New("key pFieldTypeId is not valid id.")
		}

		if formField.PFormFieldName == "" {
			return errors.New("Key pFieldName cant be empty.")
		}

		if formField.PFormFieldElement == "" {
			return errors.New("Key pFormFieldElement cant be empty.")
		}
	}

	// START: Transaction
	tx, errBegin := p.db.Begin()
	if errBegin != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			return errRollback
		}
		log.Printf("errBegin: %v", errBegin)
		return errBegin
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
			return errRollback
		}
		log.Printf("errScanQueryPartner: %v", errScanQueryPartner)
		return errScanQueryPartner
	}

	if countPartner != 1 {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollbacl: %v", errRollback)
			return errRollback
		}
		return errors.New("partner not found.")
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
			return errRollback
		}
		log.Println("errExecQueryInsertForm:", errExecQueryInsertForm)
		return errExecQueryInsertForm
	}

	rowsAffectedQueryInsertForm, errRowsAffectedQueryInsertForm := resultQueryInsertPForm.RowsAffected()
	if errRowsAffectedQueryInsertForm != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			return errRollback
		}
		return errRowsAffectedQueryInsertForm
	}
	if rowsAffectedQueryInsertForm != int64(1) {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			return errRollback
		}
		return errors.New("failed insert pForm")
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
			return errRollback
		}
		return errScanQueryCheckFieldTypeByIds
	}

	if errQueryCheckFieldTypeByIds := rowQueryCheckFieldTypeByIds.Err(); errQueryCheckFieldTypeByIds != nil {
		return errQueryCheckFieldTypeByIds
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
			return errRollback
		}
		log.Println("errQueryInserFormFields:", errQueryInserFormFields)
		return errQueryInserFormFields
	}

	rowsAffectedQueryInsertFormFields, errRowsAffectedQueryInserFormFields := resultQueryInserFormFields.RowsAffected()
	if errRowsAffectedQueryInserFormFields != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errRollback
		}
		log.Println("errRowsAffectedQueryInserFormFields:", errRowsAffectedQueryInserFormFields)
		return errRowsAffectedQueryInserFormFields
	}

	if lenFormFields != int(rowsAffectedQueryInsertFormFields) {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errRollback
		}
		return errors.New("failed insert pFormField")
	}

	if errCommit := tx.Commit(); errCommit != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errRollback
		}
		return errCommit
	}
	// END: Transaction

	return nil
}

func (p *pFormFieldUsecase) GetFormFieldById(ctx context.Context, id string) (payload.PFormFieldResponse, error) {
	return payload.PFormFieldResponse{}, nil
}
