package dto

type CreateContest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type UpdateContest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}
