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

type postgresFilesRepo struct {
	db *pgxpool.Pool
}

func NewFilesRepository(pool *pgxpool.Pool) Files {
	return &postgresFilesRepo{db: pool}
}

const createFileQuery = `INSERT INTO files(user_id, type, content_type, name, size, status, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

func (p *postgresFilesRepo) Create(ctx context.Context, file *domain.File) error {
	var id domain.ID
	createdAt := time.Now()

	row := p.db.QueryRow(ctx, createFileQuery, file.UserID, file.Type, file.ContentType, file.Name, file.Size,
		file.Status)
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}

	file.ID = id
	file.CreatedAt = createdAt

	return nil
}

const updateFileStatusQuery = `UPDATE files SET status = $1, updated_at = now() WHERE id = $2`

func (p *postgresFilesRepo) UpdateStatus(ctx context.Context, fileID domain.ID, status domain.FileStatus) error {
	tag, err := p.db.Exec(ctx, updateFileStatusQuery, status, fileID)
	if err != nil {
		return fmt.Errorf("cannot update file status: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return domain.NotFound
	}
	return nil
}

const getFileForUploading = `UPDATE files SET status = $2, updated_at = now() 
		WHERE id = (SELECT id FROM files WHERE status = $1 LIMIT 1) 
		RETURNING id, user_id, type, content_type, name, size, status, created_at, updated_at`

func (p *postgresFilesRepo) GetForUploading(ctx context.Context) (domain.File, error) {
	var file domain.File
	row := p.db.QueryRow(ctx, getFileForUploading, domain.UploadedByClient, domain.StorageUploadInProgress)
	if err := row.Scan(&file.ID, &file.UserID, &file.Type, &file.ContentType, &file.Name, &file.Size, &file.Status,
		&file.CreatedAt, &file.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.File{}, domain.NotFound
		}
		return domain.File{}, fmt.Errorf("cannot get file for uploading: %w", err)
	}

	return file, nil
}

const updateFileStatusAndURLQuery = `UPDATE files SET status = $1, url = $2, updated_at = now() WHERE id = $3`

func (p *postgresFilesRepo) UpdateStatusAndSetURL(ctx context.Context, fileID domain.ID, url string) error {
	tag, err := p.db.Exec(ctx, updateFileStatusAndURLQuery, domain.UploadedToStorage, url, fileID)
	if err != nil {
		return fmt.Errorf("cannot update status and url of file: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return domain.NotFound
	}

	return nil
}

const getFileByID = `SELECT user_id, type, content_type, name, size, status, url, created_at, updated_at FROM files 
		WHERE id = $1`

func (p *postgresFilesRepo) GetByID(ctx context.Context, fileID domain.ID) (domain.File, error) {
	var file domain.File

	row := p.db.QueryRow(ctx, getFileByID, fileID)
	if err := row.Scan(&file.UserID, &file.Type, &file.ContentType, &file.Name, &file.Size, &file.Status, &file.URL,
		&file.CreatedAt, &file.UpdatedAt); err != nil {
		return domain.File{}, fmt.Errorf("cannot get file by id: %w", err)
	}

	file.ID = fileID
	return file, nil
}
