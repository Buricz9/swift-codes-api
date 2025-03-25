package service

import (
	"database/sql"
	"fmt"
	"strings"

	"swift-codes-api/internal/repository"
)

type SwiftService interface {
	GetSwiftCodeWithBranches(code string) (interface{}, error)
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

type SwiftCodeResponseHQ struct {
	SwiftCode     string           `json:"swiftCode"`
	BankName      string           `json:"bankName"`
	Address       string           `json:"address"`
	CountryISO2   string           `json:"countryISO2"`
	CountryName   string           `json:"countryName"`
	IsHeadquarter bool             `json:"isHeadquarter"`
	Branches      []SwiftCodeBasic `json:"branches,omitempty"`
}

type SwiftCodeResponseBR struct {
	SwiftCode     string `json:"swiftCode"`
	BankName      string `json:"bankName"`
	Address       string `json:"address"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
}

type SwiftCodeBasic struct {
	SwiftCode     string `json:"swiftCode"`
	BankName      string `json:"bankName"`
	Address       string `json:"address"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
}

type CountrySwiftCodesResponse struct {
	CountryISO2 string           `json:"countryISO2"`
	CountryName string           `json:"countryName"`
	SwiftCodes  []SwiftCodeBasic `json:"swiftCodes"`
}

type swiftService struct {
	repo repository.SwiftRepository
}

func NewSwiftService(repo repository.SwiftRepository) SwiftService {
	return &swiftService{
		repo: repo,
	}
}

func (s *swiftService) GetSwiftCodeWithBranches(code string) (interface{}, error) {
	swiftCode, err := s.repo.GetBySwiftCode(code)
	if err != nil {
		return nil, fmt.Errorf("service error getting swift code: %w", err)
	}
	if swiftCode == nil {
		return nil, fmt.Errorf("swift code not found: %s", code)
	}

	if swiftCode.IsHeadquarter {
		hqResp := SwiftCodeResponseHQ{
			SwiftCode:     swiftCode.SwiftCode,
			BankName:      swiftCode.BankName,
			Address:       swiftCode.Address,
			CountryISO2:   swiftCode.CountryISO2,
			CountryName:   swiftCode.CountryName,
			IsHeadquarter: swiftCode.IsHeadquarter,
		}
		branches, err := s.repo.GetBranchesByHeadquarterCode(swiftCode.SwiftCode)
		if err != nil {
			return nil, fmt.Errorf("service error getting branches: %w", err)
		}
		for _, branch := range branches {
			brResp := SwiftCodeBasic{
				SwiftCode:     branch.SwiftCode,
				BankName:      branch.BankName,
				Address:       branch.Address,
				CountryISO2:   branch.CountryISO2,
				IsHeadquarter: branch.IsHeadquarter,
			}
			hqResp.Branches = append(hqResp.Branches, brResp)
		}
		return &hqResp, nil
	} else {
		brResp := SwiftCodeResponseBR{
			SwiftCode:     swiftCode.SwiftCode,
			BankName:      swiftCode.BankName,
			Address:       swiftCode.Address,
			CountryISO2:   swiftCode.CountryISO2,
			CountryName:   swiftCode.CountryName,
			IsHeadquarter: swiftCode.IsHeadquarter,
		}
		return &brResp, nil
	}
}

func (s *swiftService) GetSwiftCodesByCountry(countryISO2 string) (*CountrySwiftCodesResponse, error) {
	swiftCodes, err := s.repo.GetByCountryISO2(countryISO2)
	if err != nil {
		return nil, fmt.Errorf("service error getting swift codes by country: %w", err)
	}

	if len(swiftCodes) == 0 {
		return nil, fmt.Errorf("no swift codes found for country: %s", countryISO2)
	}

	var dtos []SwiftCodeBasic
	for _, sc := range swiftCodes {
		dtos = append(dtos, SwiftCodeBasic{
			SwiftCode:     sc.SwiftCode,
			BankName:      sc.BankName,
			Address:       sc.Address,
			CountryISO2:   sc.CountryISO2,
			IsHeadquarter: sc.IsHeadquarter,
		})
	}

	return &CountrySwiftCodesResponse{
		CountryISO2: countryISO2,
		CountryName: swiftCodes[0].CountryName,
		SwiftCodes:  dtos,
	}, nil
}

func (s *swiftService) CreateSwiftCode(input CreateSwiftCodeInput) error {
	countryISO2 := strings.ToUpper(input.CountryISO2)
	countryName := strings.ToUpper(input.CountryName)

	swift := repository.SwiftCode{
		SwiftCode:     input.SwiftCode,
		BankName:      input.BankName,
		Address:       input.Address,
		CountryISO2:   countryISO2,
		CountryName:   countryName,
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
