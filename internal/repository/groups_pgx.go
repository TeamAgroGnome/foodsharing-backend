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

type postgresGroupsRepo struct {
	db *pgxpool.Pool
}

func NewGroupsRepository(pool *pgxpool.Pool) Groups {
	return &postgresGroupsRepo{db: pool}
}

const createGroupQuery = `INSERT INTO groups(name, permissions, created_at) VALUES ($1, $2, $3) RETURNING id`

func (p *postgresGroupsRepo) Create(ctx context.Context, group *domain.Group) error {
	var id domain.ID
	createdAt := time.Now()

	row := p.db.QueryRow(ctx, createGroupQuery, group.Name, group.Permissions, createdAt)
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("cannot create group: %w", err)
	}

	group.ID = id
	group.CreatedAt = createdAt
	return nil
}

const updateGroupQuery = `UPDATE groups SET name = $1, permissions = $2, updated_at = now() WHERE id = $3`

func (p *postgresGroupsRepo) Update(ctx context.Context, group domain.Group) error {
	tag, err := p.db.Exec(ctx, updateGroupQuery, group.Name, group.Permissions, group.ID)
	if err != nil {
		return fmt.Errorf("cannot update group: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return domain.NotFound
	}
	return nil
}

const deleteGroupQuery = `DELETE FROM groups WHERE id = $1`

func (p *postgresGroupsRepo) Delete(ctx context.Context, id domain.ID) error {
	_, err := p.db.Exec(ctx, deleteGroupQuery, id)
	if err != nil {
		return fmt.Errorf("cannot delete group: %w", err)
	}
	return nil
}

const getGroupByID = `SELECT name, permissions, created_at, updated_at FROM groups WHERE id = $1`

func (p *postgresGroupsRepo) GetByID(ctx context.Context, id domain.ID) (domain.Group, error) {
	var group domain.Group
	row := p.db.QueryRow(ctx, getGroupByID, id)

	if err := row.Scan(&group.Name, &group.Permissions, &group.CreatedAt, &group.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Group{}, domain.NotFound
		}

		return domain.Group{}, fmt.Errorf("cannot get user by id: %w", err)
	}

	group.ID = id
	return group, nil
}

const getGroupsByName = `SELECT id, name, permissions, created_at, updated_at FROM groups WHERE name LIKE $1`

func (p *postgresGroupsRepo) GetByName(ctx context.Context, name string) ([]domain.Group, error) {
	var groups []domain.Group
	rows, err := p.db.Query(ctx, getGroupsByName, name)
	if err != nil {
		return nil, fmt.Errorf("cannot get groups by name: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var group domain.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.Permissions, &group.CreatedAt, &group.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan group: %w", err)
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot get groups by name: %w", err)
	}

	return groups, nil
}

const getGroupsByPermissionsQuery = `SELECT id, name, permissions, created_at, updated_at FROM groups 
		WHERE permissions & $1 > 0`

func (p *postgresGroupsRepo) GetByPermissions(ctx context.Context, permission domain.Permission) ([]domain.Group, error) {
	var groups []domain.Group
	rows, err := p.db.Query(ctx, getGroupsByPermissionsQuery, permission)
	if err != nil {
		return nil, fmt.Errorf("cannot get group by permissions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var group domain.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.Permissions, &group.CreatedAt, &group.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan group: %w", err)
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot get group by permissions: %w", err)
	}

	return groups, nil
}

const getAllGroupsQuery = `SELECT id, name, permissions, created_at, updated_at FROM groups`

func (p *postgresGroupsRepo) GetAll(ctx context.Context) ([]domain.Group, error) {
	var groups []domain.Group
	rows, err := p.db.Query(ctx, getAllGroupsQuery)
	if err != nil {
		return nil, fmt.Errorf("cannot get group by permissions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var group domain.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.Permissions, &group.CreatedAt, &group.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan group: %w", err)
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot get group by permissions: %w", err)
	}

	return groups, nil
}

const addUserToGroupQuery = `INSERT INTO users_to_groups(user_id, group_id) VALUES ($1, $2)`

func (p *postgresGroupsRepo) AddUser(ctx context.Context, groupID domain.ID, userID domain.ID) error {
	tag, err := p.db.Exec(ctx, addUserToGroupQuery, userID, groupID)
	if err != nil {
		return fmt.Errorf("cannot add user to group: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return domain.NotFound
	}
	return nil
}

const getUserGroupsQuery = `SELECT id, name, permissions, created_at, updated_at FROM groups 
		WHERE id IN (SELECT group_id FROM users_to_groups WHERE  user_id = $1)`

func (p *postgresGroupsRepo) GetUserGroups(ctx context.Context, userID domain.ID) ([]domain.Group, error) {
	var groups []domain.Group
	rows, err := p.db.Query(ctx, getUserGroupsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot get user groups: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var group domain.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.Permissions, &group.CreatedAt, &group.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan group: %w", err)
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot get user groups: %w", err)
	}

	return groups, nil
}

const removeUserFromGroupQuery = `DELETE FROM users_to_groups WHERE user_id = $1 AND group_id = $2`

func (p *postgresGroupsRepo) RemoveUser(ctx context.Context, groupID domain.ID, userID domain.ID) error {
	_, err := p.db.Exec(ctx, removeUserFromGroupQuery, userID, groupID)
	if err != nil {
		return fmt.Errorf("cannot remove user from group: %w", err)
	}
	return nil
}
