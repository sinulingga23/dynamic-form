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

	CreatePForm struct {
		Id         string
		PPartnerId string
		Name       string
		CreatedAt  time.Time
	}
)
