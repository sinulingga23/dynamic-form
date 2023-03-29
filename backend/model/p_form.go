package model

import (
	"database/sql"
	"time"
)

type (
	PForm struct {
		Id         string
		PPartnerId string
		Name       string
		CreatedAt  time.Time
		UpdatedAt  sql.NullTime
	}

	PFormFieldChild struct {
		PFieldTypeId      string
		PFIeldTypeName    string
		PFormFieldId      string
		PFormFieldName    string
		PFormFieldElement string
	}
	PFormDetail struct {
		Id               string
		Name             string
		PPartnerId       string
		PPartnerName     string
		PFormFieldChilds []PFormFieldChild
		CreatedAt        time.Time
		UpdatedAt        sql.NullTime
	}

	PFormPartner struct {
		Id           string
		Name         string
		PPartnerId   string
		PPartnerName string
		CreatedAt    time.Time
		UpdatedAt    sql.NullTime
	}

	CreatePForm struct {
		Id         string
		PPartnerId string
		Name       string
		CreatedAt  time.Time
	}
)
