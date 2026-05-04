package app

type Error struct {
	Code string

	Err error
}

const (
	ErrCodeResourceNotFound      = "ResourceNotFound"
	ErrCodeResourceAlreadyExists = "ResourceAlreadyExists"
	ErrCodeInvalidState          = "InvalidState"
	ErrCodeUnauthorized          = "Unauthorized"
	ErrCodeUnknown               = "Unknown"
)

func NewError(code string, err error) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func (e Error) Error() string {
	return e.Err.Error()
}
