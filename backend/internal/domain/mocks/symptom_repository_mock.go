package mocks

import (
	"context"
	"grab-bootcamp-be-team13-2025/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

type MockSymptomRepository struct {
	mock.Mock
}

func (m *MockSymptomRepository) FindByName(ctx context.Context, name string) ([]models.Symptom, error) {
	args := m.Called(ctx, name)
	return args.Get(0).([]models.Symptom), args.Error(1)
}

func (m *MockSymptomRepository) CreateUserSymptom(ctx context.Context, userSymptom *models.UserSymptom) error {
	args := m.Called(ctx, userSymptom)
	return args.Error(0)
}

func (m *MockSymptomRepository) GetUserSymptoms(ctx context.Context, userID uint64, fromAt string) ([]models.UserSymptom, error) {
	args := m.Called(ctx, userID, fromAt)
	return args.Get(0).([]models.UserSymptom), args.Error(1)
}

func (m *MockSymptomRepository) CreateFollowup(ctx context.Context, followup *models.Followup) error {
	args := m.Called(ctx, followup)
	return args.Error(0)
}

func (m *MockSymptomRepository) GetLatestFollowup(ctx context.Context, userID uint64) (*models.Followup, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Followup), args.Error(1)
}

func (m *MockSymptomRepository) UpdateFollowup(ctx context.Context, followup *models.Followup) error {
	args := m.Called(ctx, followup)
	return args.Error(0)
}

func (m *MockSymptomRepository) CreateDiagnosis(ctx context.Context, diagnosis *models.Diagnosis) error {
	args := m.Called(ctx, diagnosis)
	return args.Error(0)
}

func (m *MockSymptomRepository) GetLatestDiagnosis(ctx context.Context, userID uint64) (*models.Diagnosis, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Diagnosis), args.Error(1)
}

func (m *MockSymptomRepository) GetAll(ctx context.Context) ([]models.Symptom, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Symptom), args.Error(1)
} 