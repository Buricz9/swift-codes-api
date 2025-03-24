package service

import (
	"database/sql"
	"fmt"
	"swift-codes-api/internal/repository"
)

type SwiftService interface {
	GetSwiftCode(code string) (*repository.SwiftCode, error)
	GetSwiftCodesByCountry(countryISO2 string) ([]repository.SwiftCode, error)
	CreateSwiftCode(input CreateSwiftCodeInput) error
	DeleteSwiftCode(code string) error
}

type CreateSwiftCodeInput struct {
	SwiftCode            string
	BankName             string
	Address              string
	CountryISO2          string
	CountryName          string
	IsHeadquarter        bool
	HeadquarterSwiftCode *string
}

type swiftService struct {
	repo repository.SwiftRepository
}

func NewSwiftService(repo repository.SwiftRepository) SwiftService {
	return &swiftService{
		repo: repo,
	}
}

func (s *swiftService) GetSwiftCode(code string) (*repository.SwiftCode, error) {
	swiftCode, err := s.repo.GetBySwiftCode(code)
	if err != nil {
		return nil, fmt.Errorf("service error getting swift code: %w", err)
	}

	if swiftCode == nil {
		return nil, fmt.Errorf("swift code not found: %s", code)
	}

	return swiftCode, nil
}

func (s *swiftService) GetSwiftCodesByCountry(countryISO2 string) ([]repository.SwiftCode, error) {
	swiftCodes, err := s.repo.GetByCountryISO2(countryISO2)
	if err != nil {
		return nil, fmt.Errorf("service error getting swift codes by country: %w", err)
	}

	if swiftCodes == nil {
		return nil, fmt.Errorf("no swift codes found for country: %s", countryISO2)
	}

	return swiftCodes, nil
}

func (s *swiftService) CreateSwiftCode(input CreateSwiftCodeInput) error {
	swift := repository.SwiftCode{
		SwiftCode:     input.SwiftCode,
		BankName:      input.BankName,
		Address:       input.Address,
		CountryISO2:   input.CountryISO2,
		CountryName:   input.CountryName,
		IsHeadquarter: input.IsHeadquarter,
	}

	if input.HeadquarterSwiftCode != nil {
		swift.HeadquarterSwiftCode.String = *input.HeadquarterSwiftCode
		swift.HeadquarterSwiftCode.Valid = true
	} else {
		swift.HeadquarterSwiftCode.Valid = false
	}

	return s.repo.CreateSwiftCode(swift)
}

func (s *swiftService) DeleteSwiftCode(code string) error {
	err := s.repo.DeleteBySwiftCode(code)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("swift code not found: %s", code)
		}
		return fmt.Errorf("service error deleting swift code: %w", err)
	}

	return nil
}
