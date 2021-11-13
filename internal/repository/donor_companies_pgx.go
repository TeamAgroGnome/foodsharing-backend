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

type postgresDonorCompaniesRepo struct {
	db *pgxpool.Pool
}

func NewDonorCompaniesRepository(pool *pgxpool.Pool) DonorCompanies {
	return &postgresDonorCompaniesRepo{db: pool}
}

const createDonorCompanyQuery = `INSERT INTO donor_companies(name, city_id, contract_date, contract_number, created_at) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`

func (p *postgresDonorCompaniesRepo) Create(ctx context.Context, company *domain.DonorCompany) error {
	var id domain.ID
	createdAt := time.Now()

	row := p.db.QueryRow(ctx, createDonorCompanyQuery, company.Name, company.CityID, company.ContractDate,
		company.ContractNumber, createdAt)
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("cannot create donor company: %w", err)
	}

	company.ID = id
	company.CreatedAt = createdAt

	return nil
}

const updateDonorCompanyQuery = `UPDATE donor_companies SET name = $1, city_id = $2, contract_date = $3, 
		contract_number = $4, updated_at = now() WHERE id = $5`

func (p *postgresDonorCompaniesRepo) Update(ctx context.Context, company domain.DonorCompany) error {
	tag, err := p.db.Exec(ctx, updateDonorCompanyQuery, company.Name, company.CityID, company.ContractDate,
		company.ContractNumber, company.ID)
	if err != nil {
		return fmt.Errorf("cannot update donor company: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return domain.NotFound
	}
	return nil
}

const deleteDonorCompanyQuery = `DELETE FROM donor_companies WHERE id = $1`

func (p *postgresDonorCompaniesRepo) Delete(ctx context.Context, id domain.ID) error {
	_, err := p.db.Exec(ctx, deleteDonorCompanyQuery, id)
	if err != nil {
		return fmt.Errorf("cannot delete donor company: %w", err)
	}
	return nil
}

const getDonorCompanyByIDQuery = `SELECT name, city_id, contract_date, contract_number, created_at, updated_at 
		FROM donor_companies WHERE id = $1`

func (p *postgresDonorCompaniesRepo) GetByID(ctx context.Context, id domain.ID) (domain.DonorCompany, error) {
	var donor domain.DonorCompany
	row := p.db.QueryRow(ctx, getDonorCompanyByIDQuery, id)
	if err := row.Scan(&donor.Name, &donor.CityID, &donor.ContractDate, &donor.ContractNumber, &donor.CreatedAt,
		&donor.UpdatedAt, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.DonorCompany{}, domain.NotFound
		}
		return domain.DonorCompany{}, fmt.Errorf("cannot get donor company by id: %w", err)
	}

	donor.ID = id
	return donor, nil
}

const getDonorCompaniesByCityQuery = `SELECT id, name, contract_date, contract_number, created_at, updated_at FROM 
		donor_companies WHERE city_id LIKE $1`

func (p *postgresDonorCompaniesRepo) GetByCity(ctx context.Context, cityID domain.ID) ([]domain.DonorCompany, error) {
	var donors []domain.DonorCompany
	rows, err := p.db.Query(ctx, getDonorCompaniesByCityQuery, cityID)
	if err != nil {
		return nil, fmt.Errorf("cannot get donor companies by city: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var donor domain.DonorCompany
		if err := rows.Scan(&donor.ID, &donor.Name, &donor.ContractDate, &donor.ContractNumber, &donor.CreatedAt,
			&donor.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan company: %w", err)
		}

		donor.CityID = cityID
		donors = append(donors, donor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot get donor companies by city: %w", err)
	}

	return donors, nil
}

const getAllDonorCompaniesQuery = `SELECT id, name, city_id, contract_date, contract_number, created_at, updated_at 
		FROM donor_companies`

func (p *postgresDonorCompaniesRepo) GetAll(ctx context.Context) ([]domain.DonorCompany, error) {
	var donors []domain.DonorCompany
	rows, err := p.db.Query(ctx, getAllDonorCompaniesQuery)
	if err != nil {
		return nil, fmt.Errorf("cannot get donor companies by city: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var donor domain.DonorCompany
		if err := rows.Scan(&donor.ID, &donor.Name, &donor.CityID, &donor.ContractDate, &donor.ContractNumber,
			&donor.CreatedAt, &donor.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan company: %w", err)
		}

		donors = append(donors, donor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot get donor companies by city: %w", err)
	}

	return donors, nil
}
