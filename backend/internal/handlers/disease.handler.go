// internal/handlers/disease.handler.go
package handlers

import (
	"context"
	"net/http"

	"grab-bootcamp-be-team13-2025/internal/domain/repository"
	"grab-bootcamp-be-team13-2025/internal/domain/models"
	"github.com/gin-gonic/gin"
)

type DiseaseHandler struct {
	diseaseRepo repository.DiseaseRepository
	cache      *DiseaseCache
}

func NewDiseaseHandler(diseaseRepo repository.DiseaseRepository, cache *DiseaseCache) *DiseaseHandler {
	return &DiseaseHandler{diseaseRepo: diseaseRepo, cache: cache}
}

// GetDiseaseDescription handles GET /disease?disease_name=...
func (h *DiseaseHandler) GetDiseaseDescription(c *gin.Context) {
	diseaseName := c.Query("disease_name")
	if diseaseName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "disease_name is required"})
		return
	}

	// Try cache first
	desc, ok := h.cache.GetDescription(diseaseName)
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"disease_name": diseaseName,
			"description":  desc,
		})
		return
	}

	// Fallback: query DB
	disease, err := h.diseaseRepo.FindByName(context.Background(), diseaseName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if disease == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Disease not found"})
		return
	}

	// Update cache
	h.cache.Load([]models.Disease{*disease})

	c.JSON(http.StatusOK, gin.H{
		"disease_name": disease.Name,
		"description":  disease.Description,
	})
}
