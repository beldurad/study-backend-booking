package http

import (
	"net/http"

	app "github.com/internships-backend/test-backend-beldurad/internal"
)

func SendResponseByError(err error, w http.ResponseWriter, r *http.Request) {
	appError, ok := err.(app.Error)
	if !ok {
		WriteResponseInternalError(w, r)
		return
	}
	switch appError.Code {
	case app.ErrCodeInvalidState, app.ErrCodeResourceAlreadyExists:
		writeResponseBadRequest(w, r)
	case app.ErrCodeResourceNotFound:
		WriteResponseNotFound(w, r)
	default:
		WriteResponseInternalError(w, r)
	}
}
