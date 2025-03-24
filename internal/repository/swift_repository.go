package repository

import (
	"database/sql"
	"fmt"
	"log"
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
	existing, err := r.GetBySwiftCode(swift.SwiftCode)
	if err != nil {
		return fmt.Errorf("failed to check existing swift code: %w", err)
	}

	if existing == nil {
		insertQuery := `
            INSERT INTO swift.swift_codes
            (swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_swift_code)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
        `
		_, err := r.db.Exec(insertQuery,
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
		log.Printf("[Upsert] Inserted new swift_code=%s", swift.SwiftCode)
		return nil
	}

	logDifferences(*existing, swift)

	updateQuery := `
        UPDATE swift.swift_codes
        SET
            bank_name = $2,
            address = $3,
            country_iso2 = $4,
            country_name = $5,
            is_headquarter = $6,
            headquarter_swift_code = $7
        WHERE swift_code = $1
    `
	_, err = r.db.Exec(updateQuery,
		swift.SwiftCode,
		swift.BankName,
		swift.Address,
		swift.CountryISO2,
		swift.CountryName,
		swift.IsHeadquarter,
		swift.HeadquarterSwiftCode,
	)
	if err != nil {
		return fmt.Errorf("failed to update swift code: %w", err)
	}

	log.Printf("[Upsert] Updated existing swift_code=%s", swift.SwiftCode)
	return nil
}

func logDifferences(existing SwiftCode, updated SwiftCode) {
	if existing.BankName != updated.BankName {
		log.Printf("[Upsert] SwiftCode=%s: BankName changed from %q to %q",
			existing.SwiftCode, existing.BankName, updated.BankName)
	}
	if existing.Address != updated.Address {
		log.Printf("[Upsert] SwiftCode=%s: Address changed from %q to %q",
			existing.SwiftCode, existing.Address, updated.Address)
	}
	if existing.CountryISO2 != updated.CountryISO2 {
		log.Printf("[Upsert] SwiftCode=%s: CountryISO2 changed from %q to %q",
			existing.SwiftCode, existing.CountryISO2, updated.CountryISO2)
	}
	if existing.CountryName != updated.CountryName {
		log.Printf("[Upsert] SwiftCode=%s: CountryName changed from %q to %q",
			existing.SwiftCode, existing.CountryName, updated.CountryName)
	}
	if existing.IsHeadquarter != updated.IsHeadquarter {
		log.Printf("[Upsert] SwiftCode=%s: IsHeadquarter changed from %v to %v",
			existing.SwiftCode, existing.IsHeadquarter, updated.IsHeadquarter)
	}

	oldHQ := ""
	newHQ := ""
	if existing.HeadquarterSwiftCode.Valid {
		oldHQ = existing.HeadquarterSwiftCode.String
	}
	if updated.HeadquarterSwiftCode.Valid {
		newHQ = updated.HeadquarterSwiftCode.String
	}
	if oldHQ != newHQ {
		log.Printf("[Upsert] SwiftCode=%s: HeadquarterSwiftCode changed from %q to %q",
			existing.SwiftCode, oldHQ, newHQ)
	}
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
