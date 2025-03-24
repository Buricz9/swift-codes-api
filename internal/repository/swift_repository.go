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
