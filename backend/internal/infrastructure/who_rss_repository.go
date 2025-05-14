package infrastructure

import (
	"encoding/json"
	"fmt"
	models "grab-bootcamp-be-team13-2025/internal/domain/models"
	"io/ioutil"
)

type WhoRssRepository struct{}

func (r *WhoRssRepository) FetchWHOOutbreaks() ([]models.Outbreak, error) {
	// Đọc outbreak từ file mock_outbreaks.json (demo)
	body, err := ioutil.ReadFile("./internal/infrastructure/mock_outbreaks.json")
	if err != nil {
		return nil, fmt.Errorf("read mock_outbreaks.json failed: %w", err)
	}
	var outbreaks []models.Outbreak
	err = json.Unmarshal(body, &outbreaks)
	if err != nil {
		return nil, fmt.Errorf("parse mock_outbreaks.json failed: %w", err)
	}
	return outbreaks, nil
}
