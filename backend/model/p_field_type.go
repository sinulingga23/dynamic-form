package model

import (
	"database/sql"
	"time"
)

type (
	PFieldType struct {
		Id        string
		Name      string
		CreatedAt time.Time
		UpdatedAt sql.NullTime
	}

	CreatePFieldType struct {
		Id        string
		Name      string
		CreatedAt time.Time
	}
)
