package repository

import (
	"database/sql"
	"fmt"
)

type SwiftCode struct {
	ID                   int
	SwiftCode            string
	BankName             string
	Address              string
	CountryISO2          string
	CountryName          string
	IsHeadquarter        bool
	HeadquarterSwiftCode sql.NullString
}

type SwiftRepository interface {
	GetBySwiftCode(code string) (*SwiftCode, error)
	GetByCountryISO2(countryISO2 string) ([]SwiftCode, error)
	GetBranchesByHeadquarterCode(hqCode string) ([]SwiftCode, error)
	CreateSwiftCode(swift SwiftCode) error
	DeleteBySwiftCode(code string) error
}

type swiftRepository struct {
	db *sql.DB
}

func NewSwiftRepository(db *sql.DB) SwiftRepository {
	return &swiftRepository{db: db}
}

func (r *swiftRepository) GetBySwiftCode(code string) (*SwiftCode, error) {
	query := `
        SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_swift_code
        FROM swift.swift_codes
        WHERE swift_code = $1
    `
	row := r.db.QueryRow(query, code)

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

func (r *swiftRepository) GetByCountryISO2(countryISO2 string) ([]SwiftCode, error) {
	query := `
        SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_swift_code
        FROM swift.swift_codes
        WHERE country_iso2 = $1
    `
	rows, err := r.db.Query(query, countryISO2)
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

func (r *swiftRepository) GetBranchesByHeadquarterCode(hqCode string) ([]SwiftCode, error) {
	query := `
        SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_swift_code
        FROM swift.swift_codes
        WHERE headquarter_swift_code = $1
    `

	rows, err := r.db.Query(query, hqCode)
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

func (r *swiftRepository) CreateSwiftCode(swift SwiftCode) error {
	query := `
        INSERT INTO swift.swift_codes
        (swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_swift_code)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.db.Exec(query,
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

	return nil
}

func (r *swiftRepository) DeleteBySwiftCode(code string) error {
	query := `
        DELETE FROM swift.swift_codes
        WHERE swift_code = $1
    `
	res, err := r.db.Exec(query, code)
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
