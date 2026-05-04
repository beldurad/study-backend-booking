package http

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/go-chi/render"
	app "github.com/internships-backend/test-backend-beldurad/internal"
	appcontext "github.com/internships-backend/test-backend-beldurad/internal/context"
)

type dummyLoginRequest struct {
	app.Role `json:"role"`
}

type response struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// Handlers
func DummyLogin(log *slog.Logger, authService app.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dummyLoginRequest

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed dummy login")

			return
		}

		auth, err := authService.DummyAuthenticate(r.Context(), req.Role)
		if err != nil {
			SendResponseByError(err, w, r)
			return
		}

		render.JSON(w, r, response{
			AccessToken:  auth.AccessToken,
			RefreshToken: auth.RefreshToken,
		})
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(log *slog.Logger, authService app.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body")
			SendResponseByError(err, w, r)
			return
		}

		auth, err := authService.AuthenticateByEmail(r.Context(), app.Email(req.Email), req.Password)

		if err != nil {
			SendResponseByError(err, w, r)
		}

		render.JSON(w, r, response{
			AccessToken:  auth.AccessToken,
			RefreshToken: auth.RefreshToken,
		})
	}
}

type registerRequest struct {
	Email   string `json:"email"`
	Pasword string `json:"password"`
	Role    string `json:"role"`
}

func Register(log *slog.Logger, authService app.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req registerRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			SendResponseByError(err, w, r)
			return
		}
		user := app.CreateUser(
			app.Email(req.Email),
			app.Role(req.Role),
			req.Pasword,
		)
		auth, err := authService.Register(r.Context(), user)
		if err != nil {
			SendResponseByError(err, w, r)
			return
		}

		render.JSON(w, r, response{
			AccessToken:  auth.AccessToken,
			RefreshToken: auth.RefreshToken,
		})

	}
}

// ----------------------------------------------
// Middlewares

func TokenAuthMiddleware(log *slog.Logger, authService app.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				SendResponseByError(app.NewError(app.ErrCodeUnauthorized, nil), w, r)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				SendResponseByError(app.NewError(app.ErrCodeUnauthorized, nil), w, r)
				return
			}

			tokenString := parts[1]

			auth, err := authService.AuthenticateByToken(r.Context(), tokenString)

			if err != nil {
				SendResponseByError(app.NewError(app.ErrCodeUnauthorized, err), w, r)
				return
			}

			ctx := context.WithValue(r.Context(), appcontext.AuthContextKey, auth)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRoleMiddleware(log *slog.Logger, authService app.AuthService, roles []app.Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			rawAuth := r.Context().Value(appcontext.AuthContextKey)
			if rawAuth == nil {
				WriteResponseUnauthorized(w, r)
				return
			}
			auth, ok := rawAuth.(*app.Auth)
			if !ok {
				WriteResponseUnauthorized(w, r)
				return
			}
			if !slices.Contains(roles, auth.User.Role) {
				WriteResponseForbidden(w, r)
				return
			}
			next.ServeHTTP(w, r)

		})
	}
}
