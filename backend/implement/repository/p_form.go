package repository

import (
	"context"
	"database/sql"

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
