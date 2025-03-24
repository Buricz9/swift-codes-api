package service

import (
	"fmt"
	"swift-codes-api/internal/repository"
)

type SwiftService interface {
	GetSwiftCode(code string) (*repository.SwiftCode, error)
	GetSwiftCodesByCountry(countryISO2 string) ([]repository.SwiftCode, error)
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
