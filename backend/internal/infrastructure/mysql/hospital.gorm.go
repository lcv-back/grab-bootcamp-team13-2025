package mysql

import (
    "grab-bootcamp-be-team13-2025/internal/domain/models"
    "grab-bootcamp-be-team13-2025/internal/domain/repository"
    "gorm.io/gorm"
)

type HospitalGormRepository struct {
    db *gorm.DB
}

func NewHospitalGormRepository(db *gorm.DB) repository.HospitalRepository {
    return &HospitalGormRepository{db: db}
}

func (r *HospitalGormRepository) FindHospitalsByDisease(diseaseName string) ([]models.Hospital, error) {
    var hospitals []models.Hospital
    err := r.db.Joins("JOIN hospital_diseases ON hospital_diseases.hospital_id = hospitals.id").
        Joins("JOIN diseases ON diseases.id = hospital_diseases.disease_id").
        Where("diseases.name = ?", diseaseName).
        Preload("Diseases").
        Find(&hospitals).Error
    return hospitals, err
}
