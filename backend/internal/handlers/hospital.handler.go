package handlers

import (
	"net/http"
	"fmt"
	"encoding/json"
	"strings"

	"grab-bootcamp-be-team13-2025/internal/domain/repository"
	"grab-bootcamp-be-team13-2025/internal/utils"

	"github.com/gin-gonic/gin"
)

type HospitalHandler struct {
	hospitalRepo repository.HospitalRepository
}

func NewHospitalHandler(hospitalRepo repository.HospitalRepository) *HospitalHandler {
	return &HospitalHandler{hospitalRepo: hospitalRepo}
}

// FindNearestHospitals nhận vào disease_name, latitude, longitude, trả về danh sách tên bệnh viện gần nhất
func (h *HospitalHandler) FindNearestHospitals(c *gin.Context) {
	type Req struct {
		DiseaseName string  `json:"disease_name" binding:"required"`
		Latitude    float64 `json:"latitude" binding:"required"`
		Longitude   float64 `json:"longitude" binding:"required"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid params"})
		return
	}

	hospitals, err := h.hospitalRepo.FindHospitalsByDisease(req.DiseaseName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if len(hospitals) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No hospitals found for this disease"})
		return
	}

	// Tính khoảng cách và sắp xếp
	type hospitalDistance struct {
		Name     string
		Address  string
		Distance float64
	}
	var results []hospitalDistance
	for _, hos := range hospitals {
		dist := utils.Haversine(req.Latitude, req.Longitude, hos.Latitude, hos.Longitude)
		results = append(results, hospitalDistance{
			Name:     hos.Name,
			Address:  hos.Address,
			Distance: dist,
		})
	}
	// Sắp xếp theo Distance
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Distance > results[j].Distance {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
	// Trả về danh sách tên bệnh viện
	var hospitalNames []string
	for _, r := range results {
		hospitalNames = append(hospitalNames, r.Name)
	}
	c.JSON(http.StatusOK, gin.H{"hospitals": hospitalNames})
}

// ====== OpenStreetMap + Overpass API: Tìm bệnh viện gần nhất ======

type ExternalNearbyHospitalRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Radius    int     `json:"radius"` // mét, tuỳ chọn
}

type ExternalNearbyHospitalResponse struct {
	Name     string  `json:"name"`
	Address  string  `json:"address"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Distance float64 `json:"distance_km"`
}

// Handler tìm bệnh viện gần nhất dùng OSM + Overpass API
func FindNearbyHospitalsExternal() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ExternalNearbyHospitalRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: latitude and longitude required"})
			return
		}

		radius := req.Radius
		if radius <= 0 {
			radius = 10000 // default 10km
		}

		query := fmt.Sprintf(`[out:json];node(around:%d,%f,%f)[amenity=hospital];out;`, radius, req.Latitude, req.Longitude)
		url := "https://overpass-api.de/api/interpreter?data=" + query

		resp, err := http.Get(url)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to call Overpass API"})
			return
		}
		defer resp.Body.Close()

		var overpassResp struct {
			Elements []struct {
				ID   int64   `json:"id"`
				Lat  float64 `json:"lat"`
				Lon  float64 `json:"lon"`
				Tags struct {
					Name    string `json:"name"`
					Address string `json:"address"`
					Street  string `json:"addr:street"`
					City    string `json:"addr:city"`
				} `json:"tags"`
			} `json:"elements"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&overpassResp); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to parse Overpass API response"})
			return
		}

		var hospitals []ExternalNearbyHospitalResponse
		for _, el := range overpassResp.Elements {
			address := el.Tags.Address
			if address == "" {
				if el.Tags.Street != "" || el.Tags.City != "" {
					address = el.Tags.Street + ", " + el.Tags.City
				}
			}
			// Chỉ nhận các bệnh viện có đủ name, address, lat, lng
			if el.Tags.Name != "" && address != "" && el.Lat != 0 && el.Lon != 0 {
				// Bỏ qua các bệnh viện có tên chứa 'skin clinic' (không phân biệt hoa thường)
				nameLower := strings.ToLower(el.Tags.Name)
				if strings.Contains(nameLower, "skin clinic") {
					continue
				}
				hospitals = append(hospitals, ExternalNearbyHospitalResponse{
					Name:     el.Tags.Name,
					Address:  address,
					Lat:      el.Lat,
					Lng:      el.Lon,
					Distance: utils.Haversine(req.Latitude, req.Longitude, el.Lat, el.Lon),
				})
			}
		}
		// Sắp xếp theo distance_km tăng dần
		for i := 0; i < len(hospitals); i++ {
			for j := i + 1; j < len(hospitals); j++ {
				if hospitals[i].Distance > hospitals[j].Distance {
					hospitals[i], hospitals[j] = hospitals[j], hospitals[i]
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{"hospitals": hospitals})
	}
}

// Hàm tính khoảng cách Haversine
