package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sinulingga23/dynamic-form/backend/api/repository"
	"github.com/sinulingga23/dynamic-form/backend/model"
)

type pFieldTypeRepositoryImpl struct {
	db *sql.DB
}

func NewPFieldTypeRepositoryImpl(db *sql.DB) repository.IPFieldTypeRepository {
	return &pFieldTypeRepositoryImpl{db: db}
}

func (p *pFieldTypeRepositoryImpl) FindOne(ctx context.Context, id string) (model.PFieldType, error) {
	query := `
	select
		id, name, created_at, updated_at
	from
		partner.p_field_type
	where
		id = $1
	`

	row := p.db.QueryRow(query, id)

	pFieldType := model.PFieldType{}
	errScan := row.Scan(
		&pFieldType.Id,
		&pFieldType.Name,
		&pFieldType.CreatedAt,
		&pFieldType.UpdatedAt,
	)
	if errScan != nil {
		return model.PFieldType{}, errScan
	}

	if err := row.Err(); err != nil {
		return model.PFieldType{}, err
	}

	return pFieldType, nil
}

func (p *pFieldTypeRepositoryImpl) FindPFieldTypesByIds(ctx context.Context, ids []string) ([]model.PFieldType, error) {
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
		id, name, created_at, updated_at
	from
		partner.p_field_type
	where
		id in $1
	`
	rows, errQuery := p.db.Query(query, paramIds)
	if errQuery != nil {
		return []model.PFieldType{}, errQuery
	}

	pFieldTypes := make([]model.PFieldType, 0)
	for rows.Next() {
		pFieldType := model.PFieldType{}
		errScan := rows.Scan(
			&pFieldType.Id,
			&pFieldType.Name,
			&pFieldType.CreatedAt,
			&pFieldType.UpdatedAt,
		)
		if errScan != nil {
			return []model.PFieldType{}, errScan
		}

		pFieldTypes = append(pFieldTypes, pFieldType)
	}

	if err := rows.Err(); err != nil {
		return []model.PFieldType{}, err
	}

	return pFieldTypes, nil
}

func (p *pFieldTypeRepositoryImpl) Create(ctx context.Context, createPFieldType model.CreatePFieldType) error {
	query := `
	insert into partner.p_field_type
		(id, name, created_at)
	values
		($1, $2, $3)
	`

	_, errExec := p.db.Exec(query,
		createPFieldType.Id,
		createPFieldType.Name,
		createPFieldType.CreatedAt)
	if errExec != nil {
		return errExec
	}

	return nil
}

func (p *pFieldTypeRepositoryImpl) Creates(ctx context.Context, createPFieldTypes []model.CreatePFieldType) error {
	tx, errBegin := p.db.Begin()
	if errBegin != nil {
		return errBegin
	}

	paramCreates := ``
	lenCreatePFieldTypes := len(createPFieldTypes)
	for i := 0; i < lenCreatePFieldTypes; i++ {
		createPFieldType := createPFieldTypes[i]
		if i != lenCreatePFieldTypes-1 {
			paramCreates += fmt.Sprintf(`('%v','%v','%v'),`,
				createPFieldType.Id,
				createPFieldType.Name,
				createPFieldType.CreatedAt)
		} else {
			paramCreates += fmt.Sprintf(`('%v','%v','%v')`,
				createPFieldType.Id,
				createPFieldType.Name,
				createPFieldType.CreatedAt)
		}
	}

	query := `
	insert into partner.p_field_type
		(id, name, created_at)
	values
		$1
	`

	_, errExec := tx.Exec(query, paramCreates)
	if errExec != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errRollback
		}
		return errExec
	}

	if errCommit := tx.Commit(); errCommit != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errRollback
		}
		return errCommit
	}

	return nil
}
