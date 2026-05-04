package app

import (
	"context"
)

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

var (
	DummyAdminAuth = Auth{
		UserID:       DummyAdmin.ID,
		User:         DummyAdmin,
		AccessToken:  "accessAdminToken",
		RefreshToken: "refreshAdminToken",
	}
	DummyUserAuth = Auth{
		UserID:       DummyUser.ID,
		User:         DummyUser,
		AccessToken:  "accessUserToken",
		RefreshToken: "refreshUserToken",
	}
)

type Claims struct {
	UserID    string `json:"userID"`
	Role      string `json:"role"`
	TokenType string `json:"tokenType"`
}

type Auth struct {
	UserID       string
	User         *User
	AccessToken  string
	RefreshToken string
}

type AuthService interface {
	DummyAuthenticate(ctx context.Context, role Role) (*Auth, error)
	AuthenticateByEmail(ctx context.Context, email Email, password string) (*Auth, error)
	AuthenticateByToken(ctx context.Context, token string) (*Auth, error)
	Register(ctx context.Context, user *User) (*Auth, error)
}
