package usecase

import (
	"grab-bootcamp-be-team13-2025/internal/domain/models"
	"grab-bootcamp-be-team13-2025/internal/domain/repository"
	"strings"
)

type OutbreakUsecase struct {
	Repo repository.OutbreakRepository
}

func (u *OutbreakUsecase) GetWHOOutbreaksFiltered(disease, keywords string, limit, offset int) ([]models.Outbreak, error) {
	outbreaks, err := u.Repo.FetchWHOOutbreaks()
	if err != nil {
		return nil, err
	}
	// Filtering
	filtered := make([]models.Outbreak, 0)
	for _, o := range outbreaks {
		if disease != "" && !strings.Contains(strings.ToLower(o.Disease), strings.ToLower(disease)) {
			continue
		}
		if keywords != "" && !strings.Contains(strings.ToLower(o.Summary), strings.ToLower(keywords)) {
			continue
		}
		filtered = append(filtered, o)
	}
	// Pagination
	if offset < len(filtered) {
		filtered = filtered[offset:]
	}
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}
	return filtered, nil
}
