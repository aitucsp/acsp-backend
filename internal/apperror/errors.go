package apperror

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrUserNotFound           = Error("user doesn't exists")
	ErrPasswordMismatch       = Error("password mismatched")
	ErrBodyParsed             = Error("request body parsed badly")
	ErrProjectNotFound        = Error("project doesn't exists")
	ErrEmailAlreadyExists     = Error("email already exists")
	ErrProjectIsAlreadyExists = Error("this project is already exists")
	ErrCouldParseID           = Error("ID of users is not in correct format")
	ErrBadInputBody           = Error("invalid input body")
	ErrUserIDNotFound         = Error("user id not found")
	ErrParameterNotFound      = Error("parameter not found")
	ErrBadCredentials         = Error("bad connection credentials")
	ErrBadSigningMethod       = Error("invalid signing method")
	ErrBadClaimsType          = Error("token claims are not of type *tokenClaims")
)
