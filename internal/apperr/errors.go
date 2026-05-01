package apperr

type Error struct {
	Code string

	Err error
}

const (
	CodeResourceNotFound      = "ResourceNotFound"
	CodeResourceAlreadyExists = "ResourceAlreadyExists"
	CodeInvalidState          = "InvalidState"
	CodeUnauthorized          = "Unauthorized"
	CodeUnknown               = "Unknown"
)

func New(code string, err error) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func (e Error) Error() string {
	return e.Err.Error()
}
