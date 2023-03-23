package model

import (
	"database/sql"
	"time"
)

type (
	PPartner struct {
		Id          string
		Name        string
		Description string
		CreatedAt   time.Time
		UpdatedAt   sql.NullTime
	}

	CreatePPartner struct {
		Name        string
		Description string
		CreatedAt   time.Time
	}
)
