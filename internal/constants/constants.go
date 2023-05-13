package constants

const (
	UsersTable                   = "users"
	RolesTable                   = "roles"
	UserRolesTable               = "user_roles"
	MaterialsTable               = "scholar_materials"
	ArticlesTable                = "scholar_articles"
	ArticlesCommentsTable        = "scholar_article_comments"
	ArticleCommentVotesTable     = "scholar_article_comment_votes"
	CardsTable                   = "code_connection_cards"
	CardInvitationsTable         = "code_connection_invitations"
	ContestsTable                = "contests"
	CodingLabProjectsTable       = "coding_lab_projects"
	CodingLabProjectModulesTable = "coding_lab_project_modules"
	DatabaseName                 = "postgres"
)

const (
	BucketName           = "acsp-avatars"
	EndPoint             = "object.pscloud.io"
	UsersAvatarsFolder   = "user-avatars"
	ArticlesImagesFolder = "articles"
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
