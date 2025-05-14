package models

type Hospital struct {
	ID        uint64   `json:"id" gorm:"primaryKey"`
	Name      string   `json:"name"`
	Address   string   `json:"address"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Diseases  []Disease `json:"diseases" gorm:"many2many:hospital_diseases;joinForeignKey:hospital_id;joinReferences:disease_id"`
}
