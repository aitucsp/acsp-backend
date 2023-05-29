package model

type Course struct {
	ID          int            `json:"id" db:"id"`
	AuthorID    int            `json:"author_id" db:"author_id"`
	Title       string         `json:"title" db:"title"`
	Description string         `json:"description" db:"description"`
	Rating      int            `json:"rating" db:"rating"`
	ImageURL    string         `json:"image_url" db:"image_url"`
	CreatedAt   string         `json:"created_at" db:"created_at"`
	UpdatedAt   string         `json:"updated_at" db:"updated_at"`
	Modules     []CourseModule `json:"-,omitempty" db:"-,omitempty"`
}

type CourseModule struct {
	ID             int                  `json:"id" db:"id"`
	CourseID       int                  `json:"course_id" db:"course_id"`
	Title          string               `json:"title" db:"title"`
	ExpectedResult string               `json:"expected_result" db:"expected_result"`
	CreatedAt      string               `json:"created_at" db:"created_at"`
	UpdatedAt      string               `json:"updated_at" db:"updated_at"`
	Lessons        []CourseModuleLesson `json:"-,omitempty" db:"-,omitempty"`
}

type CourseModuleLesson struct {
	ID           int                         `json:"id" db:"id"`
	ModuleID     int                         `json:"module_id" db:"module_id"`
	Title        string                      `json:"title" db:"title"`
	Description  string                      `json:"description" db:"description"`
	ReferenceURL string                      `json:"reference_url,omitempty" db:"reference_url,omitempty"`
	CreatedAt    string                      `json:"created_at" db:"created_at"`
	UpdatedAt    string                      `json:"updated_at" db:"updated_at"`
	Comments     []CourseModuleLessonComment `json:"-,omitempty" db:"-,omitempty"`
}

type CourseModuleLessonComment struct {
	ID        int    `json:"id" db:"id"`
	LessonID  int    `json:"lesson_id" db:"lesson_id"`
	AuthorID  int    `json:"user_id" db:"user_id"`
	Text      string `json:"text" db:"text"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}
