package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type SwiftCode struct {
	ID                   int    `json:"-"`
	SwiftCode            string `json:"swiftCode"`
	BankName             string `json:"bankName"`
	Address              string `json:"address"`
	CountryISO2          string `json:"countryISO2"`
	CountryName          string `json:"countryName"`
	IsHeadquarter        bool   `json:"isHeadquarter"`
	HeadquarterSwiftCode sql.NullString
}

type SwiftRepository interface {
	GetBySwiftCode(ctx context.Context, code string) (*SwiftCode, error)
	GetByCountryISO2(ctx context.Context, countryISO2 string) ([]SwiftCode, error)
	GetBranchesByHeadquarterCode(ctx context.Context, hqCode string) ([]SwiftCode, error)
	CreateSwiftCode(ctx context.Context, swift SwiftCode) error
	DeleteBySwiftCode(ctx context.Context, code string) error
}

type swiftRepository struct {
	db *sql.DB
}

func NewSwiftRepository(db *sql.DB) SwiftRepository {
	return &swiftRepository{db: db}
}

func (r *swiftRepository) GetBySwiftCode(ctx context.Context, code string) (*SwiftCode, error) {
	query := `
        SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_swift_code
        FROM swift.swift_codes
        WHERE swift_code = $1
    `
	row := r.db.QueryRowContext(ctx, query, code)

	var swift SwiftCode
	err := row.Scan(
		&swift.ID,
		&swift.SwiftCode,
		&swift.BankName,
		&swift.Address,
		&swift.CountryISO2,
		&swift.CountryName,
		&swift.IsHeadquarter,
		&swift.HeadquarterSwiftCode,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get swift code: %w", err)
	}

	return &swift, nil
}

func (r *swiftRepository) GetByCountryISO2(ctx context.Context, countryISO2 string) ([]SwiftCode, error) {
	query := `
        SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_swift_code
        FROM swift.swift_codes
        WHERE country_iso2 = $1
    `
	rows, err := r.db.QueryContext(ctx, query, countryISO2)
	if err != nil {
		return nil, fmt.Errorf("failed to query swift codes by country: %w", err)
	}
	defer rows.Close()

	var swiftCodes []SwiftCode
	for rows.Next() {
		var swift SwiftCode
		err := rows.Scan(
			&swift.ID,
			&swift.SwiftCode,
			&swift.BankName,
			&swift.Address,
			&swift.CountryISO2,
			&swift.CountryName,
			&swift.IsHeadquarter,
			&swift.HeadquarterSwiftCode,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan swift code: %w", err)
		}

		swiftCodes = append(swiftCodes, swift)
	}

	return swiftCodes, nil
}

func (r *swiftRepository) GetBranchesByHeadquarterCode(ctx context.Context, hqCode string) ([]SwiftCode, error) {
	query := `
        SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_swift_code
        FROM swift.swift_codes
        WHERE headquarter_swift_code = $1
    `

	rows, err := r.db.QueryContext(ctx, query, hqCode)
	if err != nil {
		return nil, fmt.Errorf("failed to query branches: %w", err)
	}
	defer rows.Close()

	var branches []SwiftCode
	for rows.Next() {
		var swift SwiftCode
		err := rows.Scan(
			&swift.ID,
			&swift.SwiftCode,
			&swift.BankName,
			&swift.Address,
			&swift.CountryISO2,
			&swift.CountryName,
			&swift.IsHeadquarter,
			&swift.HeadquarterSwiftCode,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan branch: %w", err)
		}

		branches = append(branches, swift)
	}

	return branches, nil
}

func (r *swiftRepository) CreateSwiftCode(ctx context.Context, swift SwiftCode) error {
	existing, err := r.GetBySwiftCode(ctx, swift.SwiftCode)
	if err != nil {
		return fmt.Errorf("failed to check existing swift code: %w", err)
	}

	if existing != nil {
		log.Printf("[SKIP] Swift code %s already exists, skipping insert", swift.SwiftCode)
		return nil
	}

	const insertQuery = `
        INSERT INTO swift.swift_codes
        (swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_swift_code)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err = r.db.ExecContext(ctx, insertQuery,
		swift.SwiftCode,
		swift.BankName,
		swift.Address,
		swift.CountryISO2,
		swift.CountryName,
		swift.IsHeadquarter,
		swift.HeadquarterSwiftCode,
	)
	if err != nil {
		return fmt.Errorf("failed to insert swift code: %w", err)
	}

	log.Printf("[INSERT] New swift code %s inserted", swift.SwiftCode)
	return nil
}

func (r *swiftRepository) DeleteBySwiftCode(ctx context.Context, code string) error {
	query := `
        DELETE FROM swift.swift_codes
        WHERE swift_code = $1
    `
	res, err := r.db.ExecContext(ctx, query, code)
	if err != nil {
		return fmt.Errorf("failed to delete swift code: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
