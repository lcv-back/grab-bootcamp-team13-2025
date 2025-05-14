package repository

import "grab-bootcamp-be-team13-2025/internal/domain/models"

type HospitalRepository interface {
	FindHospitalsByDisease(diseaseName string) ([]models.Hospital, error)
}
