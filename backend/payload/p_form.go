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

	PFormField struct {
		PFieldTypeId      string `json:"pFieldTypeId"`      // data type of element
		PFormFieldName    string `json:"pFormFieldName"`    // caption
		PFormFieldElement string `json:"pFormFieldElement"` // tag
	}

	PFormRequest struct {
		PFormName   string       `json:"pFormName"`
		PPartnerId  string       `json:"pPartnerId"`
		PFormFields []PFormField `json:"pFormFields"`
	}

	FormResponse struct {
		Id         string       `json:"id"`
		PFormId    string       `json:"pFormId"`
		PFormName  string       `json:"formName"`
		FormFields []PFormField `json:"formFields"`
		CreatedAt  time.Time    `json:"createdAt"`
		UpdatedAt  time.Time    `json:"updatedAt"`
	}

	PFormFieldChild struct {
		PFieldTypeId      string `json:"pFieldTypeId"`
		PFIeldTypeName    string `json:"pFieldTypeName"`
		PFormFieldId      string `json:"pFormFieldId"`
		PFormFieldName    string `json:"pFormFieldName"`
		PFormFieldElement string `json:"pFormFieldElement"`
	}
	PFormDetailResponse struct {
		Id               string            `json:"id"`
		Name             string            `json:"name"`
		PPartnerId       string            `json:"pPartnerId"`
		PPartnerName     string            `json:"pPartnerName"`
		PFormFieldChilds []PFormFieldChild `json:"pFormFields"`
		CreatedAt        time.Time         `json:"createdAt"`
		UpdatedAt        time.Time         `json:"updatedAt"`
	}
)
