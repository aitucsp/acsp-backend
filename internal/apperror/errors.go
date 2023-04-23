package apperror

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrUserNotFound       = Error("user doesn't exists")
	ErrEmailNotFound      = Error("email not found in the database")
	ErrPasswordMismatch   = Error("password mismatched")
	ErrBodyParsed         = Error("request body parsed badly")
	ErrEmailAlreadyExists = Error("email already exists")
	ErrCouldParseID       = Error("ID of users is not in correct format")
	ErrBadInputBody       = Error("invalid input body")
	ErrUserIDNotFound     = Error("user id not found")
	ErrParameterNotFound  = Error("parameter not found")
	ErrBadCredentials     = Error("bad connection credentials")
	ErrBadSigningMethod   = Error("invalid signing method")
	ErrBadClaimsType      = Error("token claims are not of type *tokenClaims")
	ErrEnvVariableParsing = Error("environment variable not parsed")
	ErrIncorrectRole      = Error("does not have access, less permissions")
	ErrRollback           = Error("error occurred when rolling back transaction")
	ErrCreatingArticle    = Error("error occurred when creating an article")
	ErrUpdatingArticle    = Error("error occurred when updating an article")
	ErrDeletingArticle    = Error("error occurred when deleting an article")
	ErrCreatingComment    = Error("error occurred when creating a comment")
	ErrUpdatingComment    = Error("error occurred when updating a comment")
	ErrUpvoteComment      = Error("error occurred when upvoting a comment")
	ErrDownvoteComment    = Error("error occurred when downvoting a comment")
)
