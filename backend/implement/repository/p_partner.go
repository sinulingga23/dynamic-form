package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/sinulingga23/dynamic-form/backend/api/repository"
	"github.com/sinulingga23/dynamic-form/backend/model"
)

type pPartnerRepositoryImpl struct {
	db *sql.DB
}

func NewPPartnerRepositoryImpl(db *sql.DB) repository.IPPartnerRepository {
	return &pPartnerRepositoryImpl{db: db}
}

func (p *pPartnerRepositoryImpl) FindOne(ctx context.Context, id string) (model.PPartner, error) {
	query := `
	select
		id, name, description, created_at, updated_at 
	from
		partner.p_partner
	where id = $1`

	row := p.db.QueryRow(query, id)

	pPartner := model.PPartner{}
	errScan := row.Scan(
		&pPartner.Id,
		&pPartner.Name,
		&pPartner.Description,
		&pPartner.CreatedAt,
		&pPartner.UpdatedAt,
	)
	if errScan != nil {
		return model.PPartner{}, errScan
	}

	if err := row.Err(); err != nil {
		return model.PPartner{}, err
	}

	return pPartner, nil
}

func (p *pPartnerRepositoryImpl) FIndPPartnersByIds(ctx context.Context, ids []string) ([]model.PPartner, error) {
	paramIds := `(`
	lenIds := len(ids)
	for i := 0; i < lenIds; i++ {
		if i != lenIds-1 {
			paramIds += fmt.Sprintf(`'%s',`, ids[i])
		} else {
			paramIds += fmt.Sprintf(`'%s'`, ids[i])
		}
	}
	paramIds += `)`

	query := `
	select
		id, name, description, created_at, updated_at
	from
		partner.p_partner
	where id in`
	query += fmt.Sprintf(` %s`, paramIds)

	rows, errQuery := p.db.Query(query)
	if errQuery != nil {
		return []model.PPartner{}, errQuery
	}
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			log.Printf("[FIndPartnersByIds] err close rows: %v", errClose)
		}
	}()

	pPartners := make([]model.PPartner, 0)
	for rows.Next() {
		pPartner := model.PPartner{}
		errScan := rows.Scan(
			&pPartner.Id,
			&pPartner.Name,
			&pPartner.Description,
			&pPartner.CreatedAt,
			&pPartner.UpdatedAt,
		)
		if errScan != nil {
			return []model.PPartner{}, errScan
		}

		pPartners = append(pPartners, pPartner)
	}

	if err := rows.Err(); err != nil {
		return []model.PPartner{}, model.PPartnerError{Code: model.PPartnerErrorNameIsEmpty}
	}

	return pPartners, nil

}
func (p *pPartnerRepositoryImpl) Create(ctx context.Context, createPartner model.CreatePPartner) error {
	query := `
	insert into partner.p_partner
		(id, name, description, created_at)
	VALUES
		($1, $2, $3, $4)
	`

	_, errExec := p.db.Exec(query,
		createPartner.Id,
		createPartner.Name,
		createPartner.Description,
		createPartner.CreatedAt)
	if errExec != nil {
		return errExec
	}

	return nil
}

func (p *pPartnerRepositoryImpl) IsExistsById(ctx context.Context, id string) (bool, error) {

	query := `
	select
		count(id)
	from
		partner.p_partner
	where
		id = $1
	`

	row := p.db.QueryRow(query, id)

	count := 0
	errScan := row.Scan(&count)
	if errScan != nil {
		return false, errScan
	}

	if err := row.Err(); err != nil {
		return false, err
	}

	if count != 1 {
		return false, sql.ErrNoRows
	}

	return true, nil
}
