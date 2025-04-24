package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"swift-codes-api/internal/repository"

	"github.com/stretchr/testify/assert"
)

type mockSwiftRepo struct {
	GetBySwiftCodeFunc               func(ctx context.Context, code string) (*repository.SwiftCode, error)
	GetByCountryISO2Func             func(ctx context.Context, countryISO2 string) ([]repository.SwiftCode, error)
	GetBranchesByHeadquarterCodeFunc func(ctx context.Context, hqCode string) ([]repository.SwiftCode, error)
	CreateSwiftCodeFunc              func(ctx context.Context, swift repository.SwiftCode) error
	DeleteBySwiftCodeFunc            func(ctx context.Context, code string) error
}

func (m *mockSwiftRepo) GetBySwiftCode(ctx context.Context, code string) (*repository.SwiftCode, error) {
	return m.GetBySwiftCodeFunc(ctx, code)
}

func (m *mockSwiftRepo) GetByCountryISO2(ctx context.Context, countryISO2 string) ([]repository.SwiftCode, error) {
	return m.GetByCountryISO2Func(ctx, countryISO2)
}

func (m *mockSwiftRepo) GetBranchesByHeadquarterCode(ctx context.Context, hqCode string) ([]repository.SwiftCode, error) {
	return m.GetBranchesByHeadquarterCodeFunc(ctx, hqCode)
}

func (m *mockSwiftRepo) CreateSwiftCode(ctx context.Context, swift repository.SwiftCode) error {
	return m.CreateSwiftCodeFunc(ctx, swift)
}

func (m *mockSwiftRepo) DeleteBySwiftCode(ctx context.Context, code string) error {
	return m.DeleteBySwiftCodeFunc(ctx, code)
}

func TestGetSwiftCodeWithBranches_HQ(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		GetBySwiftCodeFunc: func(ctx context.Context, code string) (*repository.SwiftCode, error) {
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
		GetBranchesByHeadquarterCodeFunc: func(ctx context.Context, hqCode string) ([]repository.SwiftCode, error) {
			return []repository.SwiftCode{
				{
					ID:            2,
					SwiftCode:     "BRANCHCODE1",
					BankName:      "Headquarter Bank",
					Address:       "Branch 1 Address",
					CountryISO2:   "PL",
					CountryName:   "Poland",
					IsHeadquarter: false,
				},
				{
					ID:            3,
					SwiftCode:     "BRANCHCODE2",
					BankName:      "Headquarter Bank",
					Address:       "Branch 2 Address",
					CountryISO2:   "PL",
					CountryName:   "Poland",
					IsHeadquarter: false,
				},
			}, nil
		},
	}

	svc := NewSwiftService(mockRepo)
	ctx := context.Background()
	result, err := svc.GetSwiftCodeWithBranches(ctx, "HQCODEXXX")
	assert.NoError(t, err)
	assert.NotNil(t, result)

	hqResp, ok := result.(*SwiftCodeResponseHQ)
	assert.True(t, ok, "Expected type *SwiftCodeResponseHQ for headquarter")
	assert.True(t, hqResp.IsHeadquarter)
	assert.Equal(t, "Poland", hqResp.CountryName)
	assert.Len(t, hqResp.Branches, 2)
	assert.Equal(t, "BRANCHCODE1", hqResp.Branches[0].SwiftCode)
}

func TestGetSwiftCodeWithBranches_Branch(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		GetBySwiftCodeFunc: func(ctx context.Context, code string) (*repository.SwiftCode, error) {
			return &repository.SwiftCode{
				ID:                   4,
				SwiftCode:            "BRANCHCODEXXX",
				BankName:             "Branch Bank",
				Address:              "Branch Address",
				CountryISO2:          "PL",
				CountryName:          "Poland",
				IsHeadquarter:        false,
				HeadquarterSwiftCode: sql.NullString{String: "HQCODEXXX", Valid: true},
			}, nil
		},
		GetBranchesByHeadquarterCodeFunc: func(ctx context.Context, hqCode string) ([]repository.SwiftCode, error) {
			return nil, nil
		},
	}

	svc := NewSwiftService(mockRepo)
	ctx := context.Background()
	result, err := svc.GetSwiftCodeWithBranches(ctx, "BRANCHCODEXXX")
	assert.NoError(t, err)
	assert.NotNil(t, result)

	branchResp, ok := result.(*SwiftCodeResponseBR)
	assert.True(t, ok, "Expected type *SwiftCodeResponseBR for branch")
	assert.False(t, branchResp.IsHeadquarter)
}

func TestGetSwiftCodesByCountry_Success(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		GetByCountryISO2Func: func(ctx context.Context, countryISO2 string) ([]repository.SwiftCode, error) {
			return []repository.SwiftCode{
				{SwiftCode: "SWIFT1", CountryISO2: "PL", CountryName: "Poland"},
				{SwiftCode: "SWIFT2", CountryISO2: "PL", CountryName: "Poland"},
			}, nil
		},
	}

	svc := NewSwiftService(mockRepo)
	ctx := context.Background()
	result, err := svc.GetSwiftCodesByCountry(ctx, "PL")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "PL", result.CountryISO2)
	assert.Equal(t, "Poland", result.CountryName)
	assert.Len(t, result.SwiftCodes, 2)
}

func TestGetSwiftCodesByCountry_NotFound(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		GetByCountryISO2Func: func(ctx context.Context, countryISO2 string) ([]repository.SwiftCode, error) {
			return []repository.SwiftCode{}, nil
		},
	}

	svc := NewSwiftService(mockRepo)
	ctx := context.Background()
	result, err := svc.GetSwiftCodesByCountry(ctx, "XX")
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNoCountryCodes), "expected ErrNoCountryCodes sentinel")
}

func TestCreateSwiftCode_Success(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		GetBySwiftCodeFunc: func(ctx context.Context, code string) (*repository.SwiftCode, error) {
			return nil, nil
		},
		CreateSwiftCodeFunc: func(ctx context.Context, swift repository.SwiftCode) error {
			return nil
		},
	}

	svc := NewSwiftService(mockRepo)
	ctx := context.Background()
	input := CreateSwiftCodeInput{
		SwiftCode:     "NEWSWIFT",
		BankName:      "Test Bank",
		Address:       "Test Address",
		CountryISO2:   "PL",
		CountryName:   "Poland",
		IsHeadquarter: true,
	}
	err := svc.CreateSwiftCode(ctx, input)
	assert.NoError(t, err)
}

func TestDeleteSwiftCode_Success(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		DeleteBySwiftCodeFunc: func(ctx context.Context, code string) error {
			return nil
		},
	}
	svc := NewSwiftService(mockRepo)
	ctx := context.Background()
	err := svc.DeleteSwiftCode(ctx, "NEWSWIFT")
	assert.NoError(t, err)
}

func TestDeleteSwiftCode_NotFound(t *testing.T) {
	mockRepo := &mockSwiftRepo{
		DeleteBySwiftCodeFunc: func(ctx context.Context, code string) error {
			return sql.ErrNoRows
		},
	}
	svc := NewSwiftService(mockRepo)
	ctx := context.Background()
	err := svc.DeleteSwiftCode(ctx, "UNKNOWN")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound), "expected ErrNotFound sentinel")
}
