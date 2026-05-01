package postgres

import (
	"database/sql"
	"errors"

	"github.com/internships-backend/test-backend-beldurad/internal/apperr"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

func mapDBErr(err error) error {

	if errors.Is(err, sql.ErrNoRows) {
		return apperr.New(apperr.CodeResourceNotFound, err)
	}
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return err
	}

	if pqErr.Code == pqerror.UniqueViolation {
		return apperr.New(apperr.CodeResourceAlreadyExists, err)
	}

	if pqErr.Code == pqerror.CheckViolation {
		return apperr.New(apperr.CodeInvalidState, err)
	}

	return apperr.New(apperr.CodeUnknown, err)
}
