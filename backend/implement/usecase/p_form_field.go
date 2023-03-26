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

	if len(pFormFieldRequest.FormFields) == 0 {
		return errors.New("Key formFields cant be empty.")
	}

	lenFormFields := len(pFormFieldRequest.FormFields)
	errChan := make(chan error)
	for i := 0; i < lenFormFields; i++ {
		go func() {
			formField := pFormFieldRequest.FormFields[i]

			if formField.PFieldTypeId == "" {
				errChan <- errors.New("Key pFieldTypeId cant be empty.")
			}

			_, errParse := uuid.Parse(formField.PFieldTypeId)
			if errParse != nil {
				errChan <- errors.New("key pFieldTypeId is not valid id.")
			}

			if formField.PFormFieldName == "" {
				errChan <- errors.New("Key pFieldName cant be empty.")
			}

			if formField.PFormFieldElement == "" {
				errChan <- errors.New("Key pFormFieldElement cant be empty.")
			}
		}()
	}

	if errFromChan := <-errChan; errFromChan != nil {
		return errFromChan
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
	rowQueryPartner, errQueryPartner := tx.Query(queryCheckPartner, pFormFieldRequest.PPartnerId)
	if errQueryPartner != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			return errRollback
		}
		log.Printf("errQueryPartner: %v", errQueryPartner)
		return errQueryPartner
	}

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
		log.Printf("partner not found.")
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
		id in $1
	`
	rowQueryCheckFieldTypeByIds := tx.QueryRow(queryCheckFieldTypeByIds, paramFieldTypeIds)

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

	if lenFormFields != countFieldTypeByIds {
		return errors.New("fieldType not found.")
	}

	paramPFormFields := ``
	for i := 0; i < lenFormFields; i++ {
		formField := pFormFieldRequest.FormFields[i]
		if i != lenFormFields-1 {
			paramPFormFields += fmt.Sprintf(`(%s, %s, %s, %s, %s, %v),`,
				uuid.NewString(),
				pFormId,
				formField.PFieldTypeId,
				formField.PFormFieldName,
				formField.PFormFieldElement,
				time.Now())
		} else {
			paramPFormFields += fmt.Sprintf(`(%s, %s, %s, %s, %s, %v)`,
				uuid.NewString(),
				pFormId,
				formField.PFieldTypeId,
				formField.PFormFieldName,
				formField.PFormFieldElement,
				time.Now())
		}
	}

	queryInsertFormFields := `
	insert into partner.p_form_field
		(id, p_form_id, p_field_type_id, name, element, created_at)
	values
		$1
	`
	resultQueryInserFormFields, errQueryInserFormFields := tx.Exec(queryInsertFormFields, paramPFormFields)
	if errQueryInserFormFields != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errRollback
		}
		return errQueryInserFormFields
	}

	rowsAffectedQueryInsertFormFields, errRowsAffectedQueryInserFormFields := resultQueryInserFormFields.RowsAffected()
	if errRowsAffectedQueryInserFormFields != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errRollback
		}
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
