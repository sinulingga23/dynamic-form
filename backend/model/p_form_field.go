package model

import (
	"database/sql"
	"time"
)

type (
	PFormField struct {
		Id           string
		PFormId      string
		PFieldTypeId string
		Name         string
		CreatedAt    time.Time
		UpdatedAt    sql.NullTime
	}

	CreatePFormField struct {
		Id           string
		PFormId      string
		PFieldTypeId string
		Name         string
		CreatedAt    time.Time
	}

	PFormFieldError struct {
		Code int
	}
)

const (
	PFormFieldErrorIdEmpty = iota
	PFormFieldErrorPFormIdEmpty
	PFormFieldErrorPFieldTypeIdEmpty
	PFormFieldErrorNameEmpty
	PFormFieldErrorCreatedAtEmpty
	PFormFieldErrorUpdatedAtEmpty
)

func (p PFormFieldError) Error() string {
	switch p.Code {
	case PFormFieldErrorIdEmpty:
		return "field id cant be empty."
	case PFormFieldErrorPFormIdEmpty:
		return "field p_form_id cant be empty."
	case PFormFieldErrorPFieldTypeIdEmpty:
		return "field p_field_type_id cant be empty."
	case PFormFieldErrorNameEmpty:
		return "field name cant be empty."
	case PFormFieldErrorCreatedAtEmpty:
		return "field created_at cant be empty."
	case PFormFieldErrorUpdatedAtEmpty:
		return "field updated_at cant be empty."
	default:
		return "Unrecognized error code"
	}
}
