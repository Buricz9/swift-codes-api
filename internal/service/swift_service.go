package service

import (
	"fmt"
	"swift-codes-api/internal/repository"
)

type SwiftService interface {
	GetSwiftCode(code string) (*repository.SwiftCode, error)
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
