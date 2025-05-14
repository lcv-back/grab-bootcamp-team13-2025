package repository

import "grab-bootcamp-be-team13-2025/internal/domain/models"

type OutbreakRepository interface {
	FetchWHOOutbreaks() ([]models.Outbreak, error)
}
