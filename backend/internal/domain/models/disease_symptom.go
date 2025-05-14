package models

type DiseaseSymptom struct {
	DiseaseID uint64 `gorm:"primaryKey;column:disease_id;autoIncrement:false"`
	SymptomID uint64 `gorm:"primaryKey;column:symptom_id;autoIncrement:false"`
}

func (DiseaseSymptom) TableName() string {
	return "disease_symptoms"
}
