package apphttp

import (
	"net/http"

	"github.com/go-chi/render"
)

type HttpError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

const (
	badRequestCode    = "INVALID_REQUEST"
	badRequestMessage = "invalid request"

	unauthorizedCode    = "INVALID_REQUEST"
	unauthorizedMessage = "invalid request"

	internalErrorCode    = "INTERNAL_ERROR"
	internalErrorMessage = "internal server error"

	forbiddenCode    = "INVALID_REQUEST"
	forbiddenMessage = "invalid request"

	notFoundCode    = "INVALID_REQUEST"
	notFoundMessage = "invalid_request"
)

func writeResponse(w http.ResponseWriter, r *http.Request, code int, stringCode, message string) {
	render.Status(r, code)
	render.JSON(w, r, HttpError{
		Code:    stringCode,
		Message: message,
	})
}

func writeResponseBadRequest(w http.ResponseWriter, r *http.Request) {
	writeResponse(
		w,
		r,
		http.StatusBadRequest,
		badRequestCode,
		badRequestMessage,
	)
}

func WriteResponseUnauthorized(w http.ResponseWriter, r *http.Request) {
	writeResponse(
		w,
		r,
		http.StatusUnauthorized,
		unauthorizedCode,
		unauthorizedMessage,
	)
}

func WriteResponseInternalError(w http.ResponseWriter, r *http.Request) {
	writeResponse(
		w,
		r,
		http.StatusInternalServerError,
		internalErrorCode,
		internalErrorMessage,
	)
}

func WriteResponseForbidden(w http.ResponseWriter, r *http.Request) {
	writeResponse(
		w,
		r,
		http.StatusForbidden,
		forbiddenCode,
		forbiddenMessage,
	)
}

func WriteResponseNotFound(w http.ResponseWriter, r *http.Request) {
	writeResponse(
		w,
		r,
		http.StatusNotFound,
		notFoundCode,
		notFoundMessage,
	)
}
