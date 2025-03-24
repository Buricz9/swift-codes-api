package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"swift-codes-api/internal/repository"
	"testing"
)

type mockSwiftRepo struct {
	GetBySwiftCodeFunc   func(code string) (*repository.SwiftCode, error)
	GetByCountryISO2Func func(countryISO2 string) ([]repository.SwiftCode, error)
	CreateSwiftCodeFunc  func(swift repository.SwiftCode) error
}

func (m *mockSwiftRepo) GetBySwiftCode(code string) (*repository.SwiftCode, error) {
	return m.GetBySwiftCodeFunc(code)
}

func (m *mockSwiftRepo) GetByCountryISO2(countryISO2 string) ([]repository.SwiftCode, error) {
	return m.GetByCountryISO2Func(countryISO2)
}

func (m *mockSwiftRepo) CreateSwiftCode(swift repository.SwiftCode) error {
	return m.CreateSwiftCodeFunc(swift)
}

func TestCreateSwiftCode_Success(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		CreateSwiftCodeFunc: func(swift repository.SwiftCode) error {
			// Możesz tu zrobić dodatkowe sprawdzenia, np:
			if swift.SwiftCode == "" {
				return errors.New("swift code is required")
			}
			return nil
		},
	}

	service := NewSwiftService(mockRepo)

	input := CreateSwiftCodeInput{
		SwiftCode:     "INGBPLPWXXX",
		BankName:      "ING Bank Śląski",
		Address:       "ul. Sokolska 34, Katowice",
		CountryISO2:   "PL",
		CountryName:   "Poland",
		IsHeadquarter: true,
	}

	err := service.CreateSwiftCode(input)

	assert.NoError(t, err)
}

func TestCreateSwiftCode_RepoError(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		CreateSwiftCodeFunc: func(swift repository.SwiftCode) error {
			return errors.New("repo error")
		},
	}

	service := NewSwiftService(mockRepo)

	input := CreateSwiftCodeInput{
		SwiftCode:     "INGBPLPWXXX",
		BankName:      "ING Bank Śląski",
		Address:       "ul. Sokolska 34, Katowice",
		CountryISO2:   "PL",
		CountryName:   "Poland",
		IsHeadquarter: true,
	}

	err := service.CreateSwiftCode(input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repo error")
}

func TestGetSwiftCode_Success(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		GetBySwiftCodeFunc: func(code string) (*repository.SwiftCode, error) {
			return &repository.SwiftCode{
				ID:        1,
				SwiftCode: "DEUTDEFFXXX",
				BankName:  "Deutsche Bank",
			}, nil
		},
	}

	service := NewSwiftService(mockRepo)

	result, err := service.GetSwiftCode("DEUTDEFFXXX")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "DEUTDEFFXXX", result.SwiftCode)
}

func TestGetSwiftCode_NotFound(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		GetBySwiftCodeFunc: func(code string) (*repository.SwiftCode, error) {
			return nil, nil
		},
	}

	service := NewSwiftService(mockRepo)

	result, err := service.GetSwiftCode("UNKNOWN")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestGetSwiftCodesByCountry_Success(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		GetByCountryISO2Func: func(countryISO2 string) ([]repository.SwiftCode, error) {
			return []repository.SwiftCode{
				{
					ID:          1,
					SwiftCode:   "DEUTDEFFXXX",
					BankName:    "Deutsche Bank",
					CountryISO2: countryISO2,
				},
				{
					ID:          2,
					SwiftCode:   "COMMDEFFXXX",
					BankName:    "Commerzbank",
					CountryISO2: countryISO2,
				},
			}, nil
		},
	}

	service := NewSwiftService(mockRepo)

	result, err := service.GetSwiftCodesByCountry("DE")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "DEUTDEFFXXX", result[0].SwiftCode)
	assert.Equal(t, "COMMDEFFXXX", result[1].SwiftCode)
}

func TestGetSwiftCodesByCountry_NotFound(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		GetByCountryISO2Func: func(countryISO2 string) ([]repository.SwiftCode, error) {
			return nil, nil
		},
	}

	service := NewSwiftService(mockRepo)

	result, err := service.GetSwiftCodesByCountry("XX")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no swift codes found")
}

func TestGetSwiftCodesByCountry_RepoError(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		GetByCountryISO2Func: func(countryISO2 string) ([]repository.SwiftCode, error) {
			return nil, errors.New("repo error")
		},
	}

	service := NewSwiftService(mockRepo)

	result, err := service.GetSwiftCodesByCountry("DE")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "repo error")
}
