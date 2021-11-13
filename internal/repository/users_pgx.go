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

type postgresUsersRepo struct {
	db *pgxpool.Pool
}

func NewUsersRepository(pool *pgxpool.Pool) Users {
	return &postgresUsersRepo{db: pool}
}

const createUserQuery = `INSERT INTO users(surname, name, patronymic, date_of_birth, phone_number, email, city_id, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

func (p *postgresUsersRepo) Create(ctx context.Context, user *domain.User) error {
	var id domain.ID
	createdAt := time.Now()
	row := p.db.QueryRow(ctx, createUserQuery, user.Surname, user.Name, user.Patronymic, user.DateOfBirth,
		user.PhoneNumber, user.Email, user.CityID, createdAt)
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.NotFound
		}

		return fmt.Errorf("cannot create user: %w", err)
	}

	user.CreatedAt = createdAt
	user.ID = id
	return nil
}

const updateUserQuery = `UPDATE users SET surname = $1, name = $2, patronymic = $3, date_of_birth = $4, 
		phone_number = $5, email = $6, updated_at = now() WHERE id = $7`

func (p *postgresUsersRepo) Update(ctx context.Context, user domain.User) error {
	tag, err := p.db.Exec(ctx, updateUserQuery, user.Surname, user.Name, user.Patronymic, user.DateOfBirth,
		user.PhoneNumber, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("cannot update user: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return domain.NotFound
	}
	return nil
}

const deleteUserQuery = `DELETE FROM users WHERE id = $1`

func (p *postgresUsersRepo) Delete(ctx context.Context, id domain.ID) error {
	tag, err := p.db.Exec(ctx, deleteUserQuery, id)
	if err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return domain.NotFound
	}
	return nil
}

const getUserByIDQuery = `SELECT surname, name, patronymic, date_of_birth, phone_number, email, created_at, updated_at
		FROM users WHERE id = $1`

func (p *postgresUsersRepo) GetByID(ctx context.Context, id domain.ID) (domain.User, error) {
	var user domain.User
	row := p.db.QueryRow(ctx, getUserByIDQuery, id)
	if err := row.Scan(&user.Surname, &user.Name, &user.Patronymic,
		&user.DateOfBirth, &user.PhoneNumber, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.NotFound
		}

		return domain.User{}, fmt.Errorf("cannot get user by id: %w", err)
	}

	user.ID = id
	return user, nil
}

const getUserByEmailQuery = `SELECT id, surname, name, patronymic, date_of_birth, phone_number, created_at, updated_at 
		FROM users WHERE email = $1`

func (p *postgresUsersRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	row := p.db.QueryRow(ctx, getUserByEmailQuery, email)
	if err := row.Scan(&user.ID, &user.Surname, &user.Name, &user.Patronymic,
		&user.DateOfBirth, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.NotFound
		}

		return domain.User{}, fmt.Errorf("cannot get user by email: %w", err)
	}

	user.Email = email
	return user, nil
}

const getAllUsersQuery = `SELECT id, surname, name, patronymic, date_of_birth, phone_number, email, created_at, 
		updated_at FROM users`

func (p *postgresUsersRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	rows, err := p.db.Query(ctx, getAllUsersQuery)
	if err != nil {
		return nil, fmt.Errorf("cannot get all users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Surname, &user.Name, &user.Patronymic,
			&user.DateOfBirth, &user.PhoneNumber, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot get all users: %w", err)
	}

	return users, nil
}
