package repository

import (
	"context"
	"grab-bootcamp-be-team13-2025/internal/domain/models"
)

type DiseaseRepository interface {
	FindByID(ctx context.Context, id uint64) (*models.Disease, error)
	FindByIDs(ctx context.Context, ids []uint64) ([]models.Disease, error)
	CreateDiseaseSuggestion(ctx context.Context, suggestion *models.DiseaseSuggestion) error
	FindByName(ctx context.Context, name string) (*models.Disease, error)
}
