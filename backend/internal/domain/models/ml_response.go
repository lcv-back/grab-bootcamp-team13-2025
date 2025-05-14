package models

type PredictedDisease struct {
	Name        string  `json:"name"`
	Probability float64 `json:"probability"`
}

type MLResponse struct {
	UserID            uint64            `json:"user_id"`
	PredictedDiseases []PredictedDisease `json:"predicted_diseases"`
	SymptomFollowups  []string          `json:"symptom_followups"`
	Message           string            `json:"message"`
}

type FollowupRequest struct {
	UserID     uint64            `json:"user_id"`
	Symptoms   []string          `json:"symptoms"`
	ImagePaths []string          `json:"image_paths"`
	Answers    map[string]string `json:"answers"`
}
