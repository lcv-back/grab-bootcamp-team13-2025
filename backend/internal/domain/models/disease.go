package models

import "time"

type Disease struct {
	ID              uint64           `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"type:varchar(255);index"`
	Description     string           `json:"description"`
	DepartmentName  string           `json:"department_name"`
	RelatedSymptoms []Symptom        `json:"related_symptoms" gorm:"many2many:disease_symptoms;joinForeignKey:disease_id;joinReferences:symptom_id"`
	DiseaseSymptoms []DiseaseSymptom `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type DiseaseSuggestion struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	DiseaseID uint64    `json:"disease_id"`
	Score     float64   `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}
