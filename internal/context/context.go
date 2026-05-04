package appcontext

import (
	"context"

	app "github.com/internships-backend/test-backend-beldurad/internal"
)

type contextKey string

const (
	AuthContextKey contextKey = "user"
)

func SaveUserToContext(claims *app.Claims, c context.Context) context.Context {
	return context.WithValue(c, AuthContextKey, claims)
}

func GetUserIDFromContext(c context.Context) (string, bool) {
	auth, ok := c.Value(AuthContextKey).(*app.Auth)
	if !ok {
		return "", ok
	}
	return auth.UserID, ok
}

func GetRoleFromContext(c context.Context) (string, bool) {
	claim, ok := c.Value(AuthContextKey).(*app.Claims)
	if !ok {
		return "", ok
	}
	return claim.Role, ok
}
