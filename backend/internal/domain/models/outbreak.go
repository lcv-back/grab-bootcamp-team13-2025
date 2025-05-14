package models

type OutbreakWho struct {
	Cases       int    `json:"cases"`
	Deaths      int    `json:"deaths"`
	LastUpdated string `json:"last_updated"`
}

type Outbreak struct {
	ID      string       `json:"id"`
	Disease string       `json:"disease"`
	Summary string       `json:"summary"`
	Date    string       `json:"date"`
	Link    string       `json:"link"`
	Who     *OutbreakWho `json:"who"`
}
