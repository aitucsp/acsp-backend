package model

type Project struct {
	ID           int             `json:"id" db:"id"`
	DisciplineID int             `json:"discipline_id" db:"discipline_id"`
	Title        string          `json:"title" db:"title"`
	Description  string          `json:"description" db:"description"`
	Level        string          `json:"level" db:"level"`
	ImageURL     string          `json:"image_url" db:"image_url"`
	WorkHours    int             `json:"work_hours" db:"work_hours"`
	CreatedAt    string          `json:"created_at" db:"created_at"`
	UpdatedAt    string          `json:"updated_at" db:"updated_at"`
	Modules      []ProjectModule `json:"modules"`
}

type ProjectModule struct {
	ID           int      `json:"id" db:"id"`
	ProjectID    int      `json:"project_id" db:"project_id"`
	Title        string   `json:"title" db:"title"`
	Description  string   `json:"description" db:"description"`
	ReferenceURL string   `json:"reference_url" db:"reference_url"`
	CreatedAt    string   `json:"created_at" db:"created_at"`
	UpdatedAt    string   `json:"updated_at" db:"updated_at"`
	Project      *Project `json:"-" db:"-"`
}
