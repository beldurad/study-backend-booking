package domain

import (
	"time"

	"github.com/internships-backend/test-backend-beldurad/internal/apperr"
)

const (
	adminRole = "admin"
	userRole  = "user"
)

type User struct {
	Id        string    `db:"id"`
	Email     string    `db:"email"`
	Role      string    `db:"user_role"`
	CreatedAt time.Time `db:"created_at"`
}

func IsUserAdmin(role string) bool {
	return role == adminRole
}

func IsUserClient(role string) bool {
	return role == userRole
}

type UserService struct{}

func (*UserService) GetDummyUser(role string) (*User, error) {
	var res *User
	switch role {
	case "admin":
		res = &User{
			Id:        "1",
			Email:     "admin@test",
			Role:      "admin",
			CreatedAt: time.Now(),
		}
		return res, nil
	case "user":
		res = &User{
			Id:        "2",
			Email:     "user@test",
			Role:      "user",
			CreatedAt: time.Now(),
		}
		return res, nil
	default:
		return nil, apperr.New(apperr.CodeInvalidState, nil)
	}
}
