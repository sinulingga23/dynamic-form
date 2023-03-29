package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/sinulingga23/dynamic-form/backend/api/repository"
	"github.com/sinulingga23/dynamic-form/backend/model"
)

type pFormRepositoryImpl struct {
	db *sql.DB
}

func NewPFormRepositoryImpl(db *sql.DB) repository.IPFormRepository {
	return &pFormRepositoryImpl{db: db}
}

func (p *pFormRepositoryImpl) FindPFormsByPartnerId(ctx context.Context, partnerId string) ([]*model.PFormPartner, error) {
	query := `
	select
		pf.id, pf.name, pf.p_partner_id, pp.name, pp.created_at, pp.updated_at
	from
		partner.p_form as pf
	join
		partner.p_partner as pp
	on
		pf.p_partner_id = pp.id
	where
		pf.p_partner_id = $1
	`
	rows, errQuery := p.db.Query(query, partnerId)
	if errQuery != nil {
		return []*model.PFormPartner{}, errQuery
	}
	defer rows.Close()

	pFormsPartner := make([]*model.PFormPartner, 0)
	for rows.Next() {
		pFormPartner := model.PFormPartner{}

		errScan := rows.Scan(
			&pFormPartner.Id,
			&pFormPartner.Name,
			&pFormPartner.PPartnerId,
			&pFormPartner.PPartnerName,
			&pFormPartner.CreatedAt,
			&pFormPartner.UpdatedAt,
		)
		if errScan != nil {
			return []*model.PFormPartner{}, errScan
		}

		pFormsPartner = append(pFormsPartner, &pFormPartner)
	}

	if errr := rows.Err(); errr != nil {
		return []*model.PFormPartner{}, errr
	}

	return pFormsPartner, nil
}

func (p *pFormRepositoryImpl) FindOne(ctx context.Context, id string) (*model.PFormDetail, error) {

	// TODO: Is query shoud split into query or just one query for performance?
	query := `
	select
		pf.id,
        pf.name,
        pf.p_partner_id as p_partner_id, 
        pp.name as p_partner_name, 
        pft.id as p_field_type_id,
        pft.name as p_field_type_name,
        pff.id as p_form_field_id,
        pff.name as p_form_field_name,
    	pff.element as p_form_field_name,
        pp.created_at,
        pp.updated_at
    from
        partner.p_form as pf
    join
        partner.p_partner as pp
    on
        pf.p_partner_id = pp.id
    join
        partner.p_form_field as pff
    on
        pf.id = pff.p_form_id
    join
        partner.p_field_type as pft
    on
        pff.p_field_type_id = pft.id
    where
        pf.id = $1
	`
	rows, errQuery := p.db.Query(query, id)
	if errQuery != nil {
		return &model.PFormDetail{}, errQuery
	}
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			log.Printf("errClose: %v", errClose)
		}
	}()

	pFormDetail := model.PFormDetail{}
	pFormFieldChilds := make([]model.PFormFieldChild, 0)
	for rows.Next() {
		pFormFieldChild := model.PFormFieldChild{}

		errScan := rows.Scan(
			&pFormDetail.Id,
			&pFormDetail.Name,
			&pFormDetail.PPartnerId,
			&pFormDetail.PPartnerName,
			&pFormFieldChild.PFieldTypeId,
			&pFormFieldChild.PFIeldTypeName,
			&pFormFieldChild.PFormFieldId,
			&pFormFieldChild.PFormFieldName,
			&pFormFieldChild.PFormFieldElement,
			&pFormDetail.CreatedAt,
			&pFormDetail.UpdatedAt,
		)
		if errScan != nil {
			return &model.PFormDetail{}, errScan
		}

		pFormFieldChilds = append(pFormFieldChilds, pFormFieldChild)
	}

	if err := rows.Err(); err != nil {
		return &model.PFormDetail{}, err
	}

	pFormDetail.PFormFieldChilds = pFormFieldChilds

	return &pFormDetail, nil
}
