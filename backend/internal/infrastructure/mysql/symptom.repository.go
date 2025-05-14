package mysql

import (
	"context"
	"encoding/json"
	"grab-bootcamp-be-team13-2025/internal/domain/models"
	"grab-bootcamp-be-team13-2025/internal/domain/repository"

	"gorm.io/gorm"
)

type symptomRepository struct {
	db *gorm.DB
}

func (r *symptomRepository) CreateSymptom(ctx context.Context, symptom *models.Symptom) error {
	return r.db.WithContext(ctx).Create(symptom).Error
}

func NewSymptomRepository(db *gorm.DB) repository.SymptomRepository {
	return &symptomRepository{db: db}
}

func (r *symptomRepository) FindByName(ctx context.Context, name string) ([]models.Symptom, error) {
	var symptoms []models.Symptom
	if err := r.db.Where("name LIKE ?", "%"+name+"%").Find(&symptoms).Error; err != nil {
		return nil, err
	}
	return symptoms, nil
}

func (r *symptomRepository) CreateUserSymptom(ctx context.Context, userSymptom *models.UserSymptom) error {
	// Always marshal image paths to JSON string, even if empty
if userSymptom.ImagePaths == "" || userSymptom.ImagePaths == "[]" {
	userSymptom.ImagePaths = "[]"
} else {
	// try to unmarshal to check if it's valid, if not, marshal to JSON
	var tmp []string
	if err := json.Unmarshal([]byte(userSymptom.ImagePaths), &tmp); err != nil {
		imgJSON, _ := json.Marshal([]string{userSymptom.ImagePaths})
		userSymptom.ImagePaths = string(imgJSON)
	}
}
return r.db.Create(userSymptom).Error
}

func (r *symptomRepository) GetUserSymptoms(ctx context.Context, userID uint64, fromAt string) ([]models.UserSymptom, error) {
	var symptoms []models.UserSymptom
	query := r.db.Where("user_id = ?", userID)
	if fromAt != "" {
		query = query.Where("created_at >= ?", fromAt)
	}
	if err := query.Find(&symptoms).Error; err != nil {
		return nil, err
	}
	return symptoms, nil
}

func (r *symptomRepository) CreateFollowup(ctx context.Context, followup *models.Followup) error {
	return r.db.Create(followup).Error
}

func (r *symptomRepository) GetLatestFollowup(ctx context.Context, userID uint64) (*models.Followup, error) {
	var followup models.Followup
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").First(&followup).Error; err != nil {
		return nil, err
	}
	return &followup, nil
}

func (r *symptomRepository) UpdateFollowup(ctx context.Context, followup *models.Followup) error {
	return r.db.Save(followup).Error
}

func (r *symptomRepository) CreateDiagnosis(ctx context.Context, diagnosis *models.Diagnosis) error {
	return r.db.Create(diagnosis).Error
}

func (r *symptomRepository) GetLatestDiagnosis(ctx context.Context, userID uint64) (*models.Diagnosis, error) {
	var diagnosis models.Diagnosis
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").First(&diagnosis).Error; err != nil {
		return nil, err
	}
	return &diagnosis, nil
}

func (r *symptomRepository) GetAll(ctx context.Context) ([]models.Symptom, error) {
	var symptoms []models.Symptom
	if err := r.db.WithContext(ctx).Find(&symptoms).Error; err != nil {
		return nil, err
	}
	for i := range symptoms {
		var values []string
		if symptoms[i].PossibleValuesRaw != "" {
			if err := json.Unmarshal([]byte(symptoms[i].PossibleValuesRaw), &values); err != nil {
				return nil, err
			}
		}
		symptoms[i].PossibleValues = values
	}
	return symptoms, nil
}
