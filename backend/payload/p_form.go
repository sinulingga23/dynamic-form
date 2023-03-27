package payload

import (
	"time"
)

type (
	FormPartnerResponse struct {
		Id           string    `json:"id"`
		Name         string    `json:"name"`
		PPartnerId   string    `json:"pPartnerId"`
		PPartnerName string    `json:"pParterName"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
	}
)
