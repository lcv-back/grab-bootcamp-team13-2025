package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	models "grab-bootcamp-be-team13-2025/internal/domain/models"
	"grab-bootcamp-be-team13-2025/internal/domain/usecase"
	"grab-bootcamp-be-team13-2025/internal/infrastructure"

	"github.com/gin-gonic/gin"
)

type Rss struct {
	Channel struct {
		Items []RssItem `xml:"item"`
	} `xml:"channel"`
}

type RssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// WHO mock stats for demo (in-memory, fast lookup)
var whoStats = map[string]models.OutbreakWho{
	"measles":  {Cases: 3820, Deaths: 210, LastUpdated: "2024"},
	"cholera":  {Cases: 12000, Deaths: 350, LastUpdated: "2024"},
	"dengue":   {Cases: 55000, Deaths: 1200, LastUpdated: "2024"},
	"ebola":    {Cases: 120, Deaths: 80, LastUpdated: "2023"},
	"covid-19": {Cases: 1000000, Deaths: 50000, LastUpdated: "2024"},
	// Thêm các bệnh khác nếu muốn
}

// Lấy số liệu WHO theo tên bệnh (normalize về lower, bỏ dấu cách)
func getWhoStats(disease string) *models.OutbreakWho {
	d := strings.ToLower(strings.ReplaceAll(disease, " ", ""))
	if who, ok := whoStats[d]; ok {
		copy := who
		return &copy
	}
	return nil
}

// GetOutbreaksFromRSSHandler fetches and parses ProMED RSS feed, with cache and filtering
func GetOutbreaksFromRSSHandler(c *gin.Context) {
	diseaseQ := strings.ToLower(c.Query("disease"))
	keywordsQ := strings.ToLower(c.Query("keywords"))
	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}
	fromQ := c.Query("from")

	// Clean Architecture: gọi usecase để lấy outbreak từ WHO RSS
	// Khởi tạo repository và usecase
	repo := &infrastructure.WhoRssRepository{}
	usecase := &usecase.OutbreakUsecase{Repo: repo}

	outbreaks, err := usecase.GetWHOOutbreaksFiltered(diseaseQ, keywordsQ, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Nếu có filter from (ngày), lọc tiếp
	if fromQ != "" {
		fromTime, err := time.Parse("2006-01-02", fromQ)
		if err == nil {
			filtered := make([]models.Outbreak, 0)
			for _, o := range outbreaks {
				tParsed, err := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", o.Date)
				if err == nil && !tParsed.Before(fromTime) {
					filtered = append(filtered, o)
				}
			}
			outbreaks = filtered
		}
	}

	c.JSON(http.StatusOK, outbreaks)
}

// filterOutbreaks lọc và phân trang outbreak theo disease, keywords, limit, offset
func filterOutbreaks(data []models.Outbreak, diseaseQ, keywordsQ string, limit, offset int) []models.Outbreak {
	filtered := make([]models.Outbreak, 0)
	for _, o := range data {
		if diseaseQ != "" && !strings.Contains(strings.ToLower(o.Disease), diseaseQ) {
			continue
		}
		if keywordsQ != "" && !strings.Contains(strings.ToLower(o.Summary), keywordsQ) {
			continue
		}
		filtered = append(filtered, o)
	}
	if offset < len(filtered) {
		filtered = filtered[offset:]
	}
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}
	return filtered
}
