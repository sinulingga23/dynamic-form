package payload

import "time"

type (
	FormField struct {
		PFieldTypeId      string `json:"pFieldTypeId"`      // data type of element
		PFormFieldName    string `json:"pFormFieldName"`    // caption
		PFormFieldElement string `json:"pFormFieldElement"` // tag
	}

	PFormFieldRequest struct {
		PFormName  string      `json:"pFormName"`
		PPartnerId string      `json:"pPartnerId"`
		FormFields []FormField `json:"formFields"`
	}

	PFormFieldResponse struct {
		Id         string      `json:"id"`
		PFormId    string      `json:"pFormId"`
		PFormName  string      `json:"formName"`
		FormFields []FormField `json:"formFields"`
		CreatedAt  time.Time   `json:"createdAt"`
		UpdatedAt  time.Time   `json:"updatedAt"`
	}
)
