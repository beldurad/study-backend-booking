package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/internships-backend/test-backend-beldurad/internal/apperr"
	appcontext "github.com/internships-backend/test-backend-beldurad/internal/context"
	"github.com/internships-backend/test-backend-beldurad/internal/domain"
	apphttp "github.com/internships-backend/test-backend-beldurad/internal/http"
)

type JWTGenerator interface {
	Generate(userID string, secret string) (string, error)
}

type JWTMiddleware struct {
	secret string
	JWTGenerator
}

func NewJWTMiddleware(jwtGenerator JWTGenerator, secret string) *JWTMiddleware {
	return &JWTMiddleware{
		secret:       secret,
		JWTGenerator: jwtGenerator,
	}
}

func (j *JWTMiddleware) CheckTokenMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Error("no Authentication header in request", slog.String("request-id", middleware.GetReqID(r.Context())))
				apphttp.SendResponseByError(apperr.New(apperr.CodeUnauthorized, nil), w, r)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				log.Error("no bearer Authentication header in request", slog.String("request-id", middleware.GetReqID(r.Context())))
				apphttp.SendResponseByError(apperr.New(apperr.CodeUnauthorized, nil), w, r)
				return
			}

			tokenString := parts[1]

			claim := new(domain.Claims)
			token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(j.secret), nil
			})

			if err != nil || !token.Valid {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{
					"error": "Invalid or expired token",
				})
				return
			}

			ctx := appcontext.SaveUserToContext(claim, r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (j *JWTMiddleware) requireRole(log *slog.Logger, checkRole func(string) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := appcontext.GetRoleFromContext(r.Context())
			userId, userOk := appcontext.GetUserIDFromContext(r.Context())

			if !ok || !userOk {
				log.Error(
					"no user info in request context",
					slog.String("request-id", middleware.GetReqID(r.Context())),
				)
				apphttp.WriteResponseUnauthorized(w, r)
				return
			}
			if !checkRole(role) {
				log.Error(
					"user has no rights",
					slog.String("request-id", middleware.GetReqID(r.Context())),
					slog.String("user", userId),
					slog.String("role", role),
				)
				apphttp.WriteResponseForbidden(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (j *JWTMiddleware) RequireAdmin(log *slog.Logger) func(http.Handler) http.Handler {
	return j.requireRole(log, domain.IsUserAdmin)
}

func (j *JWTMiddleware) RequireUser(log *slog.Logger) func(http.Handler) http.Handler {
	return j.requireRole(log, domain.IsUserClient)
}

func (j *JWTMiddleware) RequireAdminOrUser(log *slog.Logger) func(http.Handler) http.Handler {
	return j.requireRole(log, func(s string) bool {
		return domain.IsUserAdmin(s) || domain.IsUserClient(s)
	})
}
