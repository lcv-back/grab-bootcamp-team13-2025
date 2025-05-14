package repository

import (
	"context"
	"grab-bootcamp-be-team13-2025/internal/domain/models"
)

// SymptomRepository defines the interface for symptom repository
type SymptomRepository interface {
	// Symptom related
	FindByName(ctx context.Context, name string) ([]models.Symptom, error)
	CreateSymptom(ctx context.Context, symptom *models.Symptom) error
	CreateUserSymptom(ctx context.Context, userSymptom *models.UserSymptom) error
	GetUserSymptoms(ctx context.Context, userID uint64, timeRange string) ([]models.UserSymptom, error)

	// Followup related
	CreateFollowup(ctx context.Context, followup *models.Followup) error
	GetLatestFollowup(ctx context.Context, userID uint64) (*models.Followup, error)
	UpdateFollowup(ctx context.Context, followup *models.Followup) error

	// Diagnosis related
	CreateDiagnosis(ctx context.Context, diagnosis *models.Diagnosis) error
	GetLatestDiagnosis(ctx context.Context, userID uint64) (*models.Diagnosis, error)
}
