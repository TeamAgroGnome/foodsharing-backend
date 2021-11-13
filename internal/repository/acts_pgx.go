package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"foodsharing-backend/internal/domain"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type postgresActsRepo struct {
	db *pgxpool.Pool
}

func NewActsRepository(pool *pgxpool.Pool) Acts {
	return &postgresActsRepo{db: pool}
}

const createActQuery = `INSERT INTO acts (user_id, donor_company_id, created_at) VALUES ($1, $2, $3) RETURNING id`

func (p *postgresActsRepo) Create(ctx context.Context, act *domain.Act) error {
	var id domain.ID
	createdAt := time.Now()
	row := p.db.QueryRow(ctx, createActQuery, act.UserID, act.DonorCompanyID, createdAt)
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("cannot create act: %w", err)
	}

	act.ID = id
	act.CreatedAt = createdAt
	return nil
}

const updateActQuery = `UPDATE acts SET donor_company_id = $1, updated_at = now() WHERE id = $2`

func (p *postgresActsRepo) Update(ctx context.Context, act domain.Act) error {
	tag, err := p.db.Exec(ctx, updateActQuery, act.DonorCompanyID, act.ID)
	if err != nil {
		return fmt.Errorf("cannot update act: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return domain.NotFound
	}
	return nil
}

const deleteActQuery = `DELETE FROM acts WHERE id = $1`

func (p *postgresActsRepo) Delete(ctx context.Context, id domain.ID) error {
	_, err := p.db.Exec(ctx, deleteActQuery, id)
	if err != nil {
		return fmt.Errorf("cannot delete act: %w", err)
	}
	return nil
}

const getActByIDQuery = `SELECT user_id, donor_company_id, created_at, updated_at FROM acts WHERE id = $1`

func (p *postgresActsRepo) GetByID(ctx context.Context, id domain.ID) (domain.Act, error) {
	var act domain.Act
	row := p.db.QueryRow(ctx, getActByIDQuery, id)
	if err := row.Scan(&act.UserID, &act.DonorCompanyID, &act.CreatedAt, &act.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Act{}, domain.NotFound
		}

		return domain.Act{}, fmt.Errorf("cannot get act: %w", err)
	}

	act.ID = id
	return act, nil
}

const getActsByUserID = `SELECT id, donor_company_id, created_at, updated_at FROM acts WHERE user_id = $1`

func (p *postgresActsRepo) GetByUserID(ctx context.Context, userID domain.ID) ([]domain.Act, error) {
	var acts []domain.Act
	rows, err := p.db.Query(ctx, getActsByUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot get acts by user id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var act domain.Act
		if err := rows.Scan(&act.ID, &act.DonorCompanyID, &act.CreatedAt, &act.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan act: %w", err)
		}
		act.UserID = userID
		acts = append(acts, act)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("cannot get acts by user id: %w", err)
	}

	return acts, nil
}

const getActsByDonorCompanyID = `SELECT id, user_id, created_at, updated_at FROM acts WHERE donor_company_id = $1`

func (p *postgresActsRepo) GetByDonorCompanyID(ctx context.Context, donorCompanyID domain.ID) ([]domain.Act, error) {
	var acts []domain.Act
	rows, err := p.db.Query(ctx, getActsByDonorCompanyID, donorCompanyID)
	if err != nil {
		return nil, fmt.Errorf("cannot get acts by donor company id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var act domain.Act
		if err := rows.Scan(&act.ID, &act.DonorCompanyID, &act.CreatedAt, &act.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan act: %w", err)
		}
		act.DonorCompanyID = donorCompanyID
		acts = append(acts, act)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("cannot get acts by donor company id: %w", err)
	}

	return acts, nil
}

const getAllActsQuery = `SELECT id, user_id, donor_company_id, created_at, updated_at FROM acts`

func (p *postgresActsRepo) GetAll(ctx context.Context) ([]domain.Act, error) {
	var acts []domain.Act
	rows, err := p.db.Query(ctx, getAllActsQuery)
	if err != nil {
		return nil, fmt.Errorf("cannot get acts by donor company id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var act domain.Act
		if err := rows.Scan(&act.ID, &act.UserID, &act.DonorCompanyID, &act.CreatedAt, &act.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan act: %w", err)
		}
		acts = append(acts, act)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("cannot get acts by donor company id: %w", err)
	}

	return acts, nil
}

const addFileToActQuery = `INSERT INTO files_to_acts(file_id, act_id) VALUES ($1, $2)`

func (p *postgresActsRepo) AddFile(ctx context.Context, actID domain.ID, fileID domain.ID) error {
	_, err := p.db.Exec(ctx, addFileToActQuery, fileID, actID)
	if err != nil {
		return fmt.Errorf("cannot add file to act: %w", err)
	}
	return nil
}

const getActFilesQuery = `SELECT id, user_id, type, content_type, name, size, status, url, created_at, updated_at 
		FROM files WHERE id IN (SELECT file_id FROM files_to_acts WHERE act_id = $1)`

func (p *postgresActsRepo) GetActFiles(ctx context.Context, actID domain.ID) ([]domain.File, error) {
	var files []domain.File

	rows, err := p.db.Query(ctx, getActFilesQuery, actID)
	if err != nil {
		return nil, fmt.Errorf("cannot get files from act: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var file domain.File
		if err := rows.Scan(&file.ID, &file.UserID, &file.Type, &file.ContentType, &file.Name, &file.Size, &file.Status,
			&file.URL, &file.CreatedAt, &file.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan file: %w", err)
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot get files from act: %w", err)
	}

	return files, nil
}

const removeFileFromActQuery = `DELETE FROM files_to_acts WHERE file_id = $1 AND act_id = $2`

func (p *postgresActsRepo) RemoveFile(ctx context.Context, fileID domain.ID, actID domain.ID) error {
	_, err := p.db.Exec(ctx, removeFileFromActQuery, fileID, actID)
	if err != nil {
		return fmt.Errorf("cannot remove file from act: %w", err)
	}
	return nil
}
