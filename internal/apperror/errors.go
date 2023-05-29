package apperror

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrUserNotFound         = Error("user doesn't exists")
	ErrEmailNotFound        = Error("email not found in the database")
	ErrPasswordMismatch     = Error("password mismatched")
	ErrBodyParsed           = Error("request body parsed badly")
	ErrEmailAlreadyExists   = Error("email already exists")
	ErrBadInputBody         = Error("invalid input body")
	ErrUserIDNotFound       = Error("user id not found")
	ErrParameterNotFound    = Error("parameter not found")
	ErrInvalidParameter     = Error("invalid parameter")
	ErrBadSigningMethod     = Error("invalid signing method")
	ErrRowsAffected         = Error("no rows affected")
	ErrBadClaimsType        = Error("token claims are not of type *tokenClaims")
	ErrEnvVariableParsing   = Error("environment variable not parsed")
	ErrIncorrectRole        = Error("does not have access, less permissions")
	ErrRollback             = Error("error occurred when rolling back transaction")
	ErrCreatingArticle      = Error("error occurred when creating an article")
	ErrUpdatingArticle      = Error("error occurred when updating an article")
	ErrDeletingArticle      = Error("error occurred when deleting an article")
	ErrCreatingComment      = Error("error occurred when creating a comment")
	ErrUpvoteComment        = Error("error occurred when upvoting a comment")
	ErrDownvoteComment      = Error("error occurred when downvoting a comment")
	ErrCreatingCard         = Error("error occurred when creating a card")
	ErrUpdatingCard         = Error("error occurred when updating a card")
	ErrDeletingCard         = Error("error occurred when deleting a card")
	ErrAnsweringCard        = Error("error occurred when answering a card")
	ErrCreatingCourse       = Error("error occurred when creating a course")
	ErrWhenUpdatingImageURL = Error("error occurred when updating image url")
	ErrNoAffectedRows       = Error("no affected rows")
	ErrCreatingDiscipline   = Error("error occurred when creating a discipline")
	ErrCreatingContest      = Error("error occurred when creating a contest")
	ErrCreatingModule       = Error("error occurred when creating a module")
)
