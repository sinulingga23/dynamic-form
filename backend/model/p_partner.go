package model

import (
	"database/sql"
	"time"
)

const (
	PPartnerErrorNameIsEmpty        = 1
	PPartnerErrorDescriptionIsEmpty = 2
)

type (
	PPartnerError struct {
		Code int
	}
	PPartner struct {
		Id          string
		Name        string
		Description string
		CreatedAt   time.Time
		UpdatedAt   sql.NullTime
	}

	CreatePPartner struct {
		Id          string
		Name        string
		Description string
		CreatedAt   time.Time
	}
)

func (p PPartnerError) Error() string {
	if p.Code == PPartnerErrorNameIsEmpty {
		return "Field name of p_partner cant be empty"
	}
	if p.Code == PPartnerErrorDescriptionIsEmpty {
		return "Field description of p_partner cant be empty"
	}

	return ""
}
