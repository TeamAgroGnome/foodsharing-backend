package repository

import (
	"context"
	"errors"
	"fmt"

	"foodsharing-backend/internal/domain"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type postgresSessionsRepo struct {
	db *pgxpool.Pool
}

func NewSessionsRepository(pool *pgxpool.Pool) Sessions {
	return &postgresSessionsRepo{db: pool}
}

const createSessionQuery = `INSERT INTO sessions (user_id, refresh_token, expires_at) VALUES ($1, $2, $3)`

func (p *postgresSessionsRepo) Create(ctx context.Context, session domain.Session) error {
	_, err := p.db.Exec(ctx, createSessionQuery, session.UserID, session.RefreshToken, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("cannot create session: %w", err)
	}
	return nil
}

const getSessionByRefreshTokenQuery = `SELECT user_id, expires_at FROM sessions 
		WHERE refresh_token = $1 AND expires_at > now()`

func (p *postgresSessionsRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.Session, error) {
	var session domain.Session
	row := p.db.QueryRow(ctx, getSessionByRefreshTokenQuery, refreshToken)
	if err := row.Scan(session.UserID, session.ExpiresAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Session{}, domain.NotFound
		}

		return domain.Session{}, fmt.Errorf("cannot get session by refresh token: %w", err)
	}

	session.RefreshToken = refreshToken
	return session, nil
}

const getSessionsByUserID = `SELECT refresh_token, expires_at FROM sessions 
		WHERE user_id = $1 AND expires_at > now()`

func (p *postgresSessionsRepo) GetByUserID(ctx context.Context, userID domain.ID) ([]domain.Session, error) {
	var sessions []domain.Session
	rows, err := p.db.Query(ctx, getSessionsByUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot get user sessions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var session domain.Session
		if err := rows.Scan(&session.RefreshToken, &session.ExpiresAt); err != nil {
			return nil, fmt.Errorf("cannot scan session: %w", err)
		}
		session.UserID = userID
		sessions = append(sessions, session)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("cannot get user sessions: %w", err)
	}

	return sessions, nil
}
