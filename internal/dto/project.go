package dto

type CreateProject struct {
	Title       string `json:"title" form:"title" db:"title"`
	Description string `json:"description" form:"description" db:"description"`
	Level       string `json:"level" form:"level" db:"level"`
	WorkHours   int    `json:"work_hours" form:"work_hours" db:"work_hours"`
}

type UpdateProject struct {
	Title       string `json:"title" form:"title" db:"title"`
	Description string `json:"description" form:"description" db:"description"`
	Level       string `json:"level" form:"level" db:"level"`
	WorkHours   int    `json:"work_hours" form:"work_hours" db:"work_hours"`
}
