package mocks

import (
	"context"
	"grab-bootcamp-be-team13-2025/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

type MockDiseaseRepository struct {
	mock.Mock
}

func (m *MockDiseaseRepository) FindByName(ctx context.Context, name string) (*models.Disease, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Disease), args.Error(1)
}

func (m *MockDiseaseRepository) FindByID(ctx context.Context, id uint64) (*models.Disease, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Disease), args.Error(1)
}

func (m *MockDiseaseRepository) FindByIDs(ctx context.Context, ids []uint64) ([]models.Disease, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]models.Disease), args.Error(1)
}

func (m *MockDiseaseRepository) CreateDiseaseSuggestion(ctx context.Context, suggestion *models.DiseaseSuggestion) error {
	args := m.Called(ctx, suggestion)
	return args.Error(0)
}

func (m *MockDiseaseRepository) GetAll(ctx context.Context) ([]models.Disease, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Disease), args.Error(1)
} 