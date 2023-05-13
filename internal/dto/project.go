package dto

type CreateProject struct {
	Title         string `json:"title" form:"title" db:"title"`
	Description   string `json:"description" form:"description" db:"description"`
	ReferenceList string `json:"reference_list" form:"reference_list" db:"reference_list"`
}

type UpdateProject struct {
	Title         string `json:"title" form:"title" db:"title"`
	Description   string `json:"description" form:"description" db:"description"`
	ReferenceList string `json:"reference_list" form:"reference_list" db:"reference_list"`
}
