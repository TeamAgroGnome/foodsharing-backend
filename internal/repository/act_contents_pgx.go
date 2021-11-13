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

type postgresActContentsRepo struct {
	db *pgxpool.Pool
}

func NewActContentsRepository(pool *pgxpool.Pool) ActContents {
	return &postgresActContentsRepo{db: pool}
}

const createActContentQuery = `INSERT INTO act_contents (act_id, number, name, count, price, expiration_date, comment, 
                          created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

func (p *postgresActContentsRepo) Create(ctx context.Context, contents ...domain.ActContent) error {
	batch := &pgx.Batch{}
	createdAt := time.Now()
	for _, content := range contents {
		batch.Queue(createActContentQuery, content.ActID, content.Number, content.Name, content.Count, content.Price,
			content.ExpirationDate, content.Comment, createdAt)
	}

	br := p.db.SendBatch(ctx, batch)
	defer br.Close()
	if _, err := br.Exec(); err != nil {
		return fmt.Errorf("cannot create act content: %w", err)
	}
	return nil
}

const updateActContentsQuery = `UPDATE act_contents SET number = $1, name = $2, count = $3, price = $4, 
		expiration_date = $5, comment = $6, updated_at = $7 WHERE id = $8`

func (p *postgresActContentsRepo) Update(ctx context.Context, contents ...domain.ActContent) error {
	batch := &pgx.Batch{}
	updatedAt := time.Now()
	for _, content := range contents {
		batch.Queue(updateActContentsQuery, content.Number, content.Name, content.Count, content.Price,
			content.ExpirationDate, content.Comment, updatedAt, content.ID)
	}

	br := p.db.SendBatch(ctx, batch)
	defer br.Close()
	if _, err := br.Exec(); err != nil {
		return fmt.Errorf("cannot update act content: %w", err)
	}
	return nil
}

const deleteActContentsQuery = `DELETE FROM act_contents WHERE id = $1`

func (p *postgresActContentsRepo) Delete(ctx context.Context, contentIDs ...domain.ID) error {
	batch := &pgx.Batch{}
	for _, id := range contentIDs {
		batch.Queue(deleteActContentsQuery, id)
	}

	br := p.db.SendBatch(ctx, batch)
	defer br.Close()
	if _, err := br.Exec(); err != nil {
		return fmt.Errorf("cannot delete act content: %w", err)
	}

	return nil
}

const getActContentByIDQuery = `SELECT act_id, number, name, count, price, expiration_date, comment FROM act_contents 
		WHERE id = $1`

func (p *postgresActContentsRepo) GetByID(ctx context.Context, id domain.ID) (domain.ActContent, error) {
	var content domain.ActContent

	row := p.db.QueryRow(ctx, getActContentByIDQuery, id)
	if err := row.Scan(&content.ActID, &content.Number, &content.Name, &content.Count, &content.Price,
		&content.ExpirationDate, &content.Comment); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ActContent{}, domain.NotFound
		}
		return domain.ActContent{}, fmt.Errorf("cannot get act contents: %w", err)
	}

	content.ID = id
	return content, nil
}

const getActContentsByActIDQuery = `SELECT id, number, name, count, price, expiration_date, comment, created_at,
       updated_at FROM act_contents WHERE act_id = $1`

func (p *postgresActContentsRepo) GetByActID(ctx context.Context, actID domain.ID) ([]domain.ActContent, error) {
	var contents []domain.ActContent

	rows, err := p.db.Query(ctx, getActContentsByActIDQuery, actID)
	if err != nil {
		return nil, fmt.Errorf("cannot get act contents: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var content domain.ActContent
		if err := rows.Scan(&content.ID, &content.Number, &content.Number, &content.Count, &content.Price,
			&content.ExpirationDate, &content.Comment, &content.CreatedAt, &content.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan act content: %w", err)
		}

		content.ActID = actID
		contents = append(contents, content)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot get act contents: %w", err)
	}

	return contents, nil
}
