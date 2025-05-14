package models

import (
	"time"
	"gorm.io/gorm"
)

type Symptom struct {
	ID                uint64   `json:"id" gorm:"primaryKey"`
	Name              string   `json:"name" gorm:"type:varchar(255);index"`
	Description       string   `json:"description"`
	PossibleValuesRaw string   `json:"possible_values_raw" gorm:"type:varchar(255)"`
	PossibleValues    []string `json:"possible_values" gorm:"type:varchar(255)"`
}

// UserSymptom represents a symptom reported by a user
type UserSymptom struct {
	ID         uint64    `json:"id" gorm:"primaryKey"`
	UserID     uint64    `json:"user_id" gorm:"index"`
	SymptomID  uint64    `json:"symptom_id" gorm:"index"`
	Name       string    `json:"name"`
	ImagePaths string    `json:"image_paths"` // Store as JSON string
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BeforeSave is a GORM hook that runs before saving
func (s *UserSymptom) BeforeSave(tx *gorm.DB) error {
	// If SymptomID is not set, try to find it from the name
	if s.SymptomID == 0 {
		var symptom Symptom
		if err := tx.Where("name = ?", s.Name).First(&symptom).Error; err == nil {
			s.SymptomID = symptom.ID
		}
	}
	return nil
}

// AfterFind is a GORM hook that runs after finding
func (s *UserSymptom) AfterFind(tx *gorm.DB) error {
	// If SymptomID is not set, try to find it from the name
	if s.SymptomID == 0 {
		var symptom Symptom
		if err := tx.Where("name = ?", s.Name).First(&symptom).Error; err == nil {
			s.SymptomID = symptom.ID
		}
	}
	return nil
}

// Followup represents a followup question for a user
type Followup struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	UserID    uint64    `json:"user_id"`
	Attempt   int       `json:"attempt"`
	Symptoms  []string  `json:"symptoms" gorm:"type:json"`
	Answers   []string  `json:"answers" gorm:"type:json"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Diagnosis represents a diagnosis for a user
type Diagnosis struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	UserID    uint64    `json:"user_id"`
	DiseaseID uint64    `json:"disease_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PredictRequest struct {
	UserID     string   `json:"user_id"`
	Symptoms   []string `json:"symptoms"`
	ImagePaths []string `json:"image_paths"`
	NumData    int      `json:"num_data"`
}