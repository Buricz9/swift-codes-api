package service

import (
	"database/sql"
	"testing"

	"swift-codes-api/internal/repository"

	"github.com/stretchr/testify/assert"
)

// Mock SwiftRepository
type mockSwiftRepo struct {
	GetBySwiftCodeFunc               func(code string) (*repository.SwiftCode, error)
	GetByCountryISO2Func             func(countryISO2 string) ([]repository.SwiftCode, error)
	GetBranchesByHeadquarterCodeFunc func(hqCode string) ([]repository.SwiftCode, error)
	CreateSwiftCodeFunc              func(swift repository.SwiftCode) error
	DeleteBySwiftCodeFunc            func(code string) error
}

func (m *mockSwiftRepo) GetBySwiftCode(code string) (*repository.SwiftCode, error) {
	return m.GetBySwiftCodeFunc(code)
}

func (m *mockSwiftRepo) GetByCountryISO2(countryISO2 string) ([]repository.SwiftCode, error) {
	return m.GetByCountryISO2Func(countryISO2)
}

func (m *mockSwiftRepo) GetBranchesByHeadquarterCode(hqCode string) ([]repository.SwiftCode, error) {
	return m.GetBranchesByHeadquarterCodeFunc(hqCode)
}

func (m *mockSwiftRepo) CreateSwiftCode(swift repository.SwiftCode) error {
	return m.CreateSwiftCodeFunc(swift)
}

func (m *mockSwiftRepo) DeleteBySwiftCode(code string) error {
	return m.DeleteBySwiftCodeFunc(code)
}

func TestGetSwiftCodeWithBranches_HQ(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		GetBySwiftCodeFunc: func(code string) (*repository.SwiftCode, error) {
			return &repository.SwiftCode{
				ID:            1,
				SwiftCode:     "HQCODEXXX",
				BankName:      "Headquarter Bank",
				Address:       "Main HQ Address",
				CountryISO2:   "PL",
				CountryName:   "Poland",
				IsHeadquarter: true,
			}, nil
		},
		GetBranchesByHeadquarterCodeFunc: func(hqCode string) ([]repository.SwiftCode, error) {
			return []repository.SwiftCode{
				{SwiftCode: "BRANCHCODE1"},
				{SwiftCode: "BRANCHCODE2"},
			}, nil
		},
	}

	svc := NewSwiftService(mockRepo)

	result, err := svc.GetSwiftCodeWithBranches("HQCODEXXX")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsHeadquarter)
	assert.Len(t, result.Branches, 2)
	assert.Equal(t, "BRANCHCODE1", result.Branches[0].SwiftCode)
}

func TestGetSwiftCodeWithBranches_Branch(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		GetBySwiftCodeFunc: func(code string) (*repository.SwiftCode, error) {
			return &repository.SwiftCode{
				ID:                   2,
				SwiftCode:            "BRANCHCODEXXX",
				BankName:             "Branch Bank",
				Address:              "Branch Address",
				CountryISO2:          "PL",
				CountryName:          "Poland",
				IsHeadquarter:        false,
				HeadquarterSwiftCode: sql.NullString{String: "HQCODEXXX", Valid: true},
			}, nil
		},
		GetBranchesByHeadquarterCodeFunc: func(hqCode string) ([]repository.SwiftCode, error) {
			return nil, nil // nie woła w tym przypadku
		},
	}

	svc := NewSwiftService(mockRepo)

	result, err := svc.GetSwiftCodeWithBranches("BRANCHCODEXXX")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsHeadquarter)
	assert.Nil(t, result.Branches)
	assert.Equal(t, "HQCODEXXX", *result.HeadquarterSwiftCode)
}

func TestGetSwiftCodesByCountry_Success(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		GetByCountryISO2Func: func(countryISO2 string) ([]repository.SwiftCode, error) {
			return []repository.SwiftCode{
				{
					SwiftCode:   "SWIFT1",
					CountryISO2: "PL",
					CountryName: "Poland",
				},
				{
					SwiftCode:   "SWIFT2",
					CountryISO2: "PL",
					CountryName: "Poland",
				},
			}, nil
		},
	}

	svc := NewSwiftService(mockRepo)

	result, err := svc.GetSwiftCodesByCountry("PL")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "PL", result.CountryISO2)
	assert.Equal(t, "Poland", result.CountryName)
	assert.Len(t, result.SwiftCodes, 2)
}

func TestGetSwiftCodesByCountry_NotFound(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		GetByCountryISO2Func: func(countryISO2 string) ([]repository.SwiftCode, error) {
			return []repository.SwiftCode{}, nil
		},
	}

	svc := NewSwiftService(mockRepo)

	result, err := svc.GetSwiftCodesByCountry("XX")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no swift codes found")
}

func TestCreateSwiftCode_Success(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		CreateSwiftCodeFunc: func(swift repository.SwiftCode) error {
			return nil
		},
	}

	svc := NewSwiftService(mockRepo)

	input := CreateSwiftCodeInput{
		SwiftCode:     "NEWSWIFT",
		BankName:      "Test Bank",
		Address:       "Test Address",
		CountryISO2:   "PL",
		CountryName:   "Poland",
		IsHeadquarter: true,
	}

	err := svc.CreateSwiftCode(input)

	assert.NoError(t, err)
}

func TestDeleteSwiftCode_Success(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		DeleteBySwiftCodeFunc: func(code string) error {
			return nil
		},
	}

	svc := NewSwiftService(mockRepo)

	err := svc.DeleteSwiftCode("NEWSWIFT")

	assert.NoError(t, err)
}

func TestDeleteSwiftCode_NotFound(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		DeleteBySwiftCodeFunc: func(code string) error {
			return sql.ErrNoRows
		},
	}

	svc := NewSwiftService(mockRepo)

	err := svc.DeleteSwiftCode("UNKNOWN")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
