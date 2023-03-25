package repository

import (
	"context"
	"database/sql"
	"fmt"
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

func (p *pFormRepositoryImpl) FindOne(ctx context.Context, id string) (model.PForm, error) {
	query := `
	select
		id, p_parner_id, name, created_at, updated_at
	from
		partner.p_form
	where
		id = $1
	`

	row := p.db.QueryRow(query, id)

	pForm := model.PForm{}
	errScan := row.Scan(
		&pForm.Id,
		&pForm.PPartnerId,
		&pForm.Name,
		&pForm.CreatedAt,
		&pForm.UpdatedAt,
	)
	if errScan != nil {
		return model.PForm{}, errScan
	}

	if err := row.Err(); err != nil {
		return model.PForm{}, err
	}

	return pForm, nil
}
func (p *pFormRepositoryImpl) FindPFormsByIds(ctx context.Context, ids []string) ([]model.PForm, error) {
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
		id, p_parner_id, name, created_at, updated_at
	from
		partner.p_form
	where
		id in $1
	`
	rows, errQuery := p.db.Query(query, paramIds)
	if errQuery != nil {
		return []model.PForm{}, errQuery
	}
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			log.Printf("[FindPFormsByIds] err close rows: %v", errClose)
		}
	}()

	pPforms := make([]model.PForm, 0)
	for rows.Next() {
		pForm := model.PForm{}
		errScan := rows.Scan(
			&pForm.Id,
			&pForm.PPartnerId,
			&pForm.Name,
			&pForm.CreatedAt,
			&pForm.UpdatedAt,
		)
		if errScan != nil {
			log.Printf("errScan: %v", errScan)
		}

		pPforms = append(pPforms, pForm)
	}

	return pPforms, nil
}
func (p *pFormRepositoryImpl) Create(ctx context.Context, createPForm model.CreatePForm) error {
	query := `
	insert into partner.p_form
		(id, p_partner_id, name, created_at)
	VALUES
		($1,$2,$3,$4)
	`

	_, errExec := p.db.Exec(query,
		createPForm.Id,
		createPForm.PPartnerId,
		createPForm.Name,
		createPForm.CreatedAt)
	if errExec != nil {
		return errExec
	}

	return nil
}
