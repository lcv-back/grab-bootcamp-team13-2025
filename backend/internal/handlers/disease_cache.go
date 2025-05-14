package handlers

import (
	"grab-bootcamp-be-team13-2025/internal/domain/models"
)

type DiseaseCache struct {
	cache map[string]string // map[disease_name]description
}

func NewDiseaseCache() *DiseaseCache {
	return &DiseaseCache{cache: make(map[string]string)}
}

func (dc *DiseaseCache) Load(diseases []models.Disease) {
	dc.cache = make(map[string]string)
	for _, d := range diseases {
		dc.cache[d.Name] = d.Description
	}
}

func (dc *DiseaseCache) GetDescription(name string) (string, bool) {
	desc, ok := dc.cache[name]
	return desc, ok
}
