package handlers

import (
	"fmt"
	models "grab-bootcamp-be-team13-2025/internal/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetOutbreaksHandler handles GET /api/outbreaks
func GetOutbreaksHandler(c *gin.Context) {
	disease := c.Query("disease")
	from := c.Query("from")
	to := c.Query("to")
	keywords := c.Query("keywords")
	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}
	// TODO: Query database or service for outbreaks with filters
	// For now, return mock data

	results := []models.Outbreak{
		{
			ID:      "12345",
			Disease: "Measles",
			Summary: "Outbreak reported in Nigeria affecting 80+ children.",
			Date:    "2025-05-13",
			Link:    "https://promedmail.org/post/xyz123",
			Who: &models.OutbreakWho{
				Cases:       3820,
				Deaths:      210,
				LastUpdated: "2024",
			},
		},
		{
			ID:      "12346",
			Disease: "Cholera",
			Summary: "WHO confirms new outbreak of cholera in Bangladesh.",
			Date:    "2025-05-12",
			Link:    "https://promedmail.org/post/xyz124",
			Who:     nil,
		},
	}

	// Filter logic can be implemented here based on query params (disease, from, to, keywords)
	_ = disease
	_ = from
	_ = to
	_ = keywords
	_ = limit
	_ = offset

	c.JSON(http.StatusOK, results)
}
