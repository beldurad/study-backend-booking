package apphttp

import (
	"net/http"

	"github.com/internships-backend/test-backend-beldurad/internal/apperr"
)

func SendResponseByError(err error, w http.ResponseWriter, r *http.Request) {
	appError, ok := err.(apperr.Error)
	if !ok {
		WriteResponseInternalError(w, r)
		return
	}
	switch appError.Code {
	case apperr.CodeInvalidState, apperr.CodeResourceAlreadyExists:
		writeResponseBadRequest(w, r)
	case apperr.CodeResourceNotFound:
		WriteResponseNotFound(w, r)
	default:
		WriteResponseInternalError(w, r)
	}
}
