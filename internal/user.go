package app

import (
	"context"
	"fmt"

	"regexp"
	"time"

	"github.com/google/uuid"
)

const (
	AdminRole = "admin"
	UserRole  = "user"
)

var (
	DummyAdmin = &User{
		ID:           "1",
		Email:        "admin@test.com",
		Role:         AdminRole,
		CreatedAt:    time.Now(),
		Password:     "1",
		PasswordHash: "1",
	}
	DummyUser = &User{
		ID:           "2",
		Email:        "user@test.com",
		Role:         UserRole,
		CreatedAt:    time.Now(),
		Password:     "1",
		PasswordHash: "1",
	}
)

type User struct {
	ID           string
	Email        Email
	Role         Role
	CreatedAt    time.Time
	Password     string
	PasswordHash string
}

func CreateUser(email Email, role Role, password string) *User {
	return &User{
		ID:        uuid.NewString(),
		Email:     email,
		Role:      role,
		Password:  password,
		CreatedAt: time.Now(),
	}
}

func (u *User) Validate() error {
	if err := uuid.Validate(u.ID); err != nil {
		return NewError(ErrCodeInvalidState, fmt.Errorf("invalid id"))
	}
	if u.Role.Validate() != nil {
		return NewError(ErrCodeInvalidState, fmt.Errorf("invalid role"))
	}
	if u.Email.Validate() != nil {
		return NewError(ErrCodeInvalidState, fmt.Errorf("invalid email"))
	}
	if u.CreatedAt.IsZero() {
		return NewError(ErrCodeInvalidState, fmt.Errorf("createdAt is required"))
	}
	if u.Password == "" {
		return NewError(ErrCodeInvalidState, fmt.Errorf("password is required"))
	}
	return nil
}

func IsUserAdmin(role string) bool {
	return role == AdminRole
}

func IsUserClient(role string) bool {
	return role == UserRole
}

type UserService interface {
	GetDummyUser(ctx context.Context, role Role) (*User, error)
	GetUserByEmail(ctx context.Context, email Email) (*User, error)
	CreateUser(ctx context.Context, user *User) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Matches(hashedPassword, password string) bool
}

const (
	emailRegex = `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$`
)

type Email string

func (e Email) Validate() error {
	valid, err := regexp.MatchString(emailRegex, string(e))
	if err != nil {
		panic(err)
	}
	if !valid {
		return NewError(ErrCodeInvalidState, nil)
	}
	return nil
}

type Role string

func (r Role) Validate() error {
	if r == AdminRole || r == UserRole {
		return nil
	}
	return NewError(ErrCodeInvalidState, nil)
}
