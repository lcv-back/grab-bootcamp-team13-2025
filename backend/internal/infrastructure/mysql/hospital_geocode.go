package mysql

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"grab-bootcamp-be-team13-2025/internal/domain/models"
	"gorm.io/gorm"
)

// GeocodeHospital cập nhật latitude/longitude cho hospital dựa trên địa chỉ bằng Nominatim API
func GeocodeHospital(db *gorm.DB, hospital *models.Hospital) error {
	if hospital.Address == "" {
		return fmt.Errorf("hospital address is empty")
	}
	endpoint := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Add("q", hospital.Address)
	params.Add("format", "json")
	params.Add("limit", "1")
	resp, err := http.Get(fmt.Sprintf("%s?%s", endpoint, params.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || len(result) == 0 {
		return fmt.Errorf("not found or decode error")
	}
	lat, err := strconv.ParseFloat(result[0].Lat, 64)
	if err != nil {
		return err
	}
	lon, err := strconv.ParseFloat(result[0].Lon, 64)
	if err != nil {
		return err
	}
	// Cập nhật vào struct và DB
	hospital.Latitude = lat
	hospital.Longitude = lon
	return db.Model(hospital).Updates(map[string]interface{}{"latitude": lat, "longitude": lon}).Error
}

// Hàm batch: cập nhật toạ độ cho tất cả hospital chưa có lat/lon
func GeocodeAllHospitals(db *gorm.DB) error {
	var hospitals []models.Hospital
	if err := db.Where("latitude = 0 OR longitude = 0").Find(&hospitals).Error; err != nil {
		return err
	}
	for _, h := range hospitals {
		err := GeocodeHospital(db, &h)
		if err != nil {
			fmt.Printf("Geocode failed for %s: %v\n", h.Name, err)
		}
	}
	return nil
}
