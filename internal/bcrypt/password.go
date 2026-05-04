package bcrypt

import (
	app "github.com/internships-backend/test-backend-beldurad/internal"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct{}

func (u *PasswordHasher) Hash(password string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", app.NewError(app.ErrCodeUnknown, err)
	}
	return string(res), nil
}

func (u *PasswordHasher) Matches(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
