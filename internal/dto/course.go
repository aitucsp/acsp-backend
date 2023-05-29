package dto

// CreateCourse DTO for Creating Course
type CreateCourse struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// UpdateCourse DTO for Updating Course
type UpdateCourse struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Rating      int    `json:"rating" validate:"required"`
}

// CreateCourseModule DTO for Creating Course Module
type CreateCourseModule struct {
	Title          string `json:"title" validate:"required"`
	ExpectedResult string `json:"expected_result" validate:"required"`
}

// UpdateCourseModule DTO for Updating Course Module
type UpdateCourseModule struct {
	Title          string `json:"title" validate:"required"`
	ExpectedResult string `json:"expected_result" validate:"required"`
}

// CreateCourseModuleLesson DTO for Creating Course Module Lesson
type CreateCourseModuleLesson struct {
	Title        string `json:"title" validate:"required"`
	Description  string `json:"description" validate:"required"`
	ReferenceURL string `json:"reference_url"`
}

// UpdateCourseModuleLesson DTO for Updating Course Module Lesson
type UpdateCourseModuleLesson struct {
	Title        string `json:"title" validate:"required"`
	Description  string `json:"description" validate:"required"`
	ReferenceURL string `json:"reference_url"`
}

// CreateLessonComment DTO for Creating Course Module Lesson Comment
type CreateLessonComment struct {
	Text string `json:"text" validate:"required"`
}

// UpdateLessonComment DTO for Updating Course Module Lesson Comment
type UpdateLessonComment struct {
	Text string `json:"text" validate:"required"`
}
