package constants

const (
	UsersTable                       = "users"
	UserDetailsTable                 = "user_details"
	RolesTable                       = "roles"
	UserRolesTable                   = "user_roles"
	MaterialsTable                   = "scholar_materials"
	ArticlesTable                    = "scholar_articles"
	ArticlesCommentsTable            = "scholar_article_comments"
	ArticleCommentVotesTable         = "scholar_article_comment_votes"
	CardsTable                       = "code_connection_cards"
	CardInvitationsTable             = "code_connection_invitations"
	ContestsTable                    = "contests"
	CodingLabDisciplinesTable        = "coding_lab_disciplines"
	CodingLabProjectsTable           = "coding_lab_projects"
	CodingLabProjectModulesTable     = "coding_lab_project_modules"
	CoursesTable                     = "courses"
	CourseModulesTable               = "course_modules"
	CourseModuleLessonsTable         = "course_module_lessons"
	CourseLessonCommentsTable        = "course_lesson_comments"
	CourseLessonCommentsAnswersTable = "course_lesson_comments_answers"
	DatabaseName                     = "postgres"
)

const (
	BucketName                  = "acsp-avatars"
	EndPoint                    = "object.pscloud.io"
	UsersAvatarsFolder          = "user-avatars"
	ArticlesImagesFolder        = "articles"
	MaterialsImagesFolder       = "materials"
	ProjectsImagesFolder        = "projects"
	ProjectsModulesImagesFolder = "projects-modules"
	DisciplinesImagesFolder     = "disciplines"
)

type VoteType int

const (
	ContextTimeoutSeconds   = 10
	FallBackDurationSeconds = 10
	DefaultUserRoleID       = 1 // ID of 'user' role
	UpvoteType              = 1
	DownvoteType            = -1
)

const (
	AcceptedStatus = "ACCEPTED"
	DeclinedStatus = "DECLINED"
)
