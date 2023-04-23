package constants

const (
	UsersTable               = "users"
	RolesTable               = "roles"
	UserRolesTable           = "user_roles"
	ArticlesTable            = "scholar_articles"
	CommentsTable            = "scholar_comments"
	MaterialsTable           = "scholar_materials"
	CommentVotesTable        = "scholar_article_comment_votes"
	CardsTable               = "code_connection_cards"
	InvitationsTable         = "code_connection_invitations"
	InvitationResponsesTable = "code_connection_responses"
	DatabaseName             = "postgres"
)

type VoteType int

const (
	ContextTimeoutSeconds   = 10
	FallBackDurationSeconds = 10
	DefaultUserRoleID       = 1 // ID of 'user' role
	UpvoteType              = 1
	DownvoteType            = -1
)
