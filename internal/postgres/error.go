package postgres

import (
	"database/sql"
	"errors"

	app "github.com/internships-backend/test-backend-beldurad/internal"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

func mapDBErr(err error) error {

	if errors.Is(err, sql.ErrNoRows) {
		return app.NewError(app.ErrCodeResourceNotFound, err)
	}
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return err
	}

	if pqErr.Code == pqerror.UniqueViolation {
		return app.NewError(app.ErrCodeResourceAlreadyExists, err)
	}

	if pqErr.Code == pqerror.CheckViolation {
		return app.NewError(app.ErrCodeInvalidState, err)
	}

	return app.NewError(app.ErrCodeUnknown, err)
}
