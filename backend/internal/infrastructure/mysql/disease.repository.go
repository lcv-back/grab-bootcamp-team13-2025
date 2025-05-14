// internal/infrastructure/mysql/disease_repository.go
package mysql

import (
	"context"
	"errors"
	"fmt"
	"grab-bootcamp-be-team13-2025/internal/domain/models"
	"grab-bootcamp-be-team13-2025/internal/domain/repository"

	"gorm.io/gorm"
)

type MySQLDiseaseRepository struct {
	db *gorm.DB
}

func NewMySQLDiseaseRepository(db *gorm.DB) repository.DiseaseRepository {
	return &MySQLDiseaseRepository{db: db}
}

func (r *MySQLDiseaseRepository) FindByID(ctx context.Context, id uint64) (*models.Disease, error) {
	var disease models.Disease
	if err := r.db.WithContext(ctx).
		Preload("RelatedSymptoms").
		Where("id = ?", id).
		First(&disease).Error; err != nil {
		return nil, err
	}
	return &disease, nil
}

func (r *MySQLDiseaseRepository) FindByIDs(ctx context.Context, ids []uint64) ([]models.Disease, error) {
	var diseases []models.Disease
	if err := r.db.WithContext(ctx).
		Preload("RelatedSymptoms").
		Where("id IN ?", ids).
		Find(&diseases).Error; err != nil {
		return nil, err
	}
	return diseases, nil
}

func (r *MySQLDiseaseRepository) CreateDiseaseSuggestion(ctx context.Context, suggestion *models.DiseaseSuggestion) error {
	return r.db.WithContext(ctx).Create(suggestion).Error
}

// FindByName implements the DiseaseRepository interface to fetch a disease by its name
func (r *MySQLDiseaseRepository) FindByName(ctx context.Context, name string) (*models.Disease, error) {
	var disease models.Disease
	if err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&disease).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record is found
		}
		return nil, fmt.Errorf("failed to find disease by name: %w", err)
	}
	return &disease, nil
}
