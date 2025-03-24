package service

import (
	"database/sql"
	"fmt"
	"swift-codes-api/internal/repository"
)

type SwiftService interface {
	GetSwiftCodeWithBranches(code string) (*SwiftCodeWithBranches, error)
	GetSwiftCodesByCountry(countryISO2 string) (*CountrySwiftCodesResponse, error)
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

type SwiftCodeWithBranches struct {
	ID                   int                    `json:"id"`
	SwiftCode            string                 `json:"swiftCode"`
	BankName             string                 `json:"bankName"`
	Address              string                 `json:"address"`
	CountryISO2          string                 `json:"countryISO2"`
	CountryName          string                 `json:"countryName"`
	IsHeadquarter        bool                   `json:"isHeadquarter"`
	HeadquarterSwiftCode *string                `json:"headquarterSwiftCode,omitempty"`
	Branches             []repository.SwiftCode `json:"branches,omitempty"`
}

type CountrySwiftCodesResponse struct {
	CountryISO2 string                 `json:"countryISO2"`
	CountryName string                 `json:"countryName"`
	SwiftCodes  []repository.SwiftCode `json:"swiftCodes"`
}

type swiftService struct {
	repo repository.SwiftRepository
}

func NewSwiftService(repo repository.SwiftRepository) SwiftService {
	return &swiftService{
		repo: repo,
	}
}

func (s *swiftService) GetSwiftCodeWithBranches(code string) (*SwiftCodeWithBranches, error) {
	swiftCode, err := s.repo.GetBySwiftCode(code)
	if err != nil {
		return nil, fmt.Errorf("service error getting swift code: %w", err)
	}

	if swiftCode == nil {
		return nil, fmt.Errorf("swift code not found: %s", code)
	}

	result := &SwiftCodeWithBranches{
		ID:            swiftCode.ID,
		SwiftCode:     swiftCode.SwiftCode,
		BankName:      swiftCode.BankName,
		Address:       swiftCode.Address,
		CountryISO2:   swiftCode.CountryISO2,
		CountryName:   swiftCode.CountryName,
		IsHeadquarter: swiftCode.IsHeadquarter,
	}

	if swiftCode.HeadquarterSwiftCode.Valid {
		hqCode := swiftCode.HeadquarterSwiftCode.String
		result.HeadquarterSwiftCode = &hqCode
	}

	if swiftCode.IsHeadquarter {
		branches, err := s.repo.GetBranchesByHeadquarterCode(swiftCode.SwiftCode)
		if err != nil {
			return nil, fmt.Errorf("service error getting branches: %w", err)
		}
		result.Branches = branches
	}

	return result, nil
}

func (s *swiftService) GetSwiftCodesByCountry(countryISO2 string) (*CountrySwiftCodesResponse, error) {
	swiftCodes, err := s.repo.GetByCountryISO2(countryISO2)
	if err != nil {
		return nil, fmt.Errorf("service error getting swift codes by country: %w", err)
	}

	if len(swiftCodes) == 0 {
		return nil, fmt.Errorf("no swift codes found for country: %s", countryISO2)
	}

	response := &CountrySwiftCodesResponse{
		CountryISO2: countryISO2,
		CountryName: swiftCodes[0].CountryName,
		SwiftCodes:  swiftCodes,
	}

	return response, nil
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
