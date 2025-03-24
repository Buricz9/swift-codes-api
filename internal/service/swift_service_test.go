package service

import (
	"errors"
	"testing"

	"swift-codes-api/internal/repository"

	"github.com/stretchr/testify/assert"
)

type mockSwiftRepo struct {
	GetBySwiftCodeFunc func(code string) (*repository.SwiftCode, error)
}

func (m *mockSwiftRepo) GetBySwiftCode(code string) (*repository.SwiftCode, error) {
	return m.GetBySwiftCodeFunc(code)
}

func TestGetSwiftCode_Success(t *testing.T) {
	// Arrange
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

	// Act
	result, err := service.GetSwiftCode("DEUTDEFFXXX")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "DEUTDEFFXXX", result.SwiftCode)
}

func TestGetSwiftCode_NotFound(t *testing.T) {
	// Arrange
	mockRepo := &mockSwiftRepo{
		GetBySwiftCodeFunc: func(code string) (*repository.SwiftCode, error) {
			return nil, nil
		},
	}

	service := NewSwiftService(mockRepo)

	// Act
	result, err := service.GetSwiftCode("UNKNOWN")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestGetSwiftCode_RepoError(t *testing.T) {
	// Arrange
	mockRepo := &mockSwiftRepo{
		GetBySwiftCodeFunc: func(code string) (*repository.SwiftCode, error) {
			return nil, errors.New("repo error")
		},
	}

	service := NewSwiftService(mockRepo)

	// Act
	result, err := service.GetSwiftCode("FAIL")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "repo error")
}
