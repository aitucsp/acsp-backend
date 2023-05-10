package dto

type CreateProject struct {
	Title         string `json:"title" db:"title"`
	Description   string `json:"description" db:"description"`
	ReferenceList string `json:"language" db:"reference_list"`
}

type UpdateProject struct {
	Title         string `json:"title" db:"title"`
	Description   string `json:"description" db:"description"`
	ReferenceList string `json:"language" db:"reference_list"`
}
