package model

type Project struct {
	ID            int             `json:"id" db:"id"`
	Title         string          `json:"title" db:"title"`
	Description   string          `json:"description" db:"description"`
	ImageURL      string          `json:"image_url" db:"image_url"`
	ReferenceList string          `json:"language" db:"reference_list"`
	CreatedAt     string          `json:"created_at" db:"created_at"`
	UpdatedAt     string          `json:"updated_at" db:"updated_at"`
	Modules       []ProjectModule `json:"modules"`
}

type ProjectModule struct {
	ID        int      `json:"id" db:"id"`
	ProjectID int      `json:"project_id" db:"project_id"`
	Title     string   `json:"title" db:"title"`
	CreatedAt string   `json:"created_at" db:"created_at"`
	UpdatedAt string   `json:"updated_at" db:"updated_at"`
	Project   *Project `json:"-" db:"-"`
}
