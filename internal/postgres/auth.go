package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	app "github.com/internships-backend/test-backend-beldurad/internal"
)

var (
	tokenSigningMethod = jwt.SigningMethodHS256
)

type AuthService struct {
	app.UserService
	app.PasswordHasher
	secret string
	db     *sql.DB
	log    *slog.Logger
}

func newAuthService(userService app.UserService, passHasher app.PasswordHasher, secret string, db *sql.DB, log *slog.Logger) *AuthService {
	return &AuthService{
		UserService:    userService,
		PasswordHasher: passHasher,
		secret:         secret,
		db:             db,
		log:            log,
	}
}

func (s *AuthService) DummyAuthenticate(ctx context.Context, role app.Role) (*app.Auth, error) {
	switch role {
	case app.AdminRole:
		return &app.DummyAdminAuth, nil
	case app.UserRole:
		return &app.DummyUserAuth, nil
	default:
		return nil, app.NewError(app.ErrCodeInvalidState, nil)
	}
}

func (s *AuthService) AuthenticateByEmail(ctx context.Context, email app.Email, password string) (*app.Auth, error) {
	if err := email.Validate(); err != nil {
		return nil, app.NewError(app.ErrCodeInvalidState, err)
	}
	user, err := s.UserService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if !s.PasswordHasher.Matches(user.PasswordHash, password) {
		return nil, app.NewError(app.ErrCodeUnauthorized, err)
	}

	auth, err := generateAuth(user, s.secret)

	tx, err := s.db.Begin()
	if err != nil {
		return nil, app.NewError(app.ErrCodeUnknown, err)
	}

	if err := saveAuth(ctx, tx, auth); err != nil {
		return nil, mapDBErr(err)
	}
	return auth, nil
}

func (s *AuthService) AuthenticateByToken(ctx context.Context, token string) (*app.Auth, error) {
	switch token {
	case app.DummyAdminAuth.AccessToken, app.DummyAdminAuth.RefreshToken:
		return &app.DummyAdminAuth, nil
	case app.DummyUserAuth.AccessToken, app.DummyUserAuth.RefreshToken:
		return &app.DummyUserAuth, nil
	}

	tokenClaims := tokenClaims{}

	jwtToken, err := jwt.ParseWithClaims(token, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})

	if err != nil || !jwtToken.Valid {
		return nil, app.NewError(app.ErrCodeUnauthorized, err)
	}

	tx, err := s.db.Begin()

	if err != nil {
		return nil, app.NewError(app.ErrCodeUnknown, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	auth, err := getAuthByUserID(ctx, tx, tokenClaims.UserID)

	if err != nil {
		return nil, mapDBErr(err)
	}

	switch tokenClaims.TokenType {
	case app.AccessTokenType:
		if auth.AccessToken != token {
			return nil, app.NewError(app.ErrCodeUnauthorized, fmt.Errorf("invalid access token"))
		}
		return auth, nil
	case app.RefreshTokenType:
		if auth.RefreshToken != token {
			return nil, app.NewError(app.ErrCodeUnauthorized, fmt.Errorf("invalid refresh token"))
		}
		auth, err := generateAuth(auth.User, s.secret)
		if err != nil {
			return nil, app.NewError(app.ErrCodeUnknown, err)
		}
		if err := updateAuth(ctx, tx, auth); err != nil {
			return nil, mapDBErr(err)
		}
		return auth, nil
	default:
		return nil, app.NewError(app.ErrCodeUnauthorized, err)
	}
}

func (s *AuthService) Register(ctx context.Context, user *app.User) (*app.Auth, error) {
	if err := user.Validate(); err != nil {
		return nil, app.NewError(app.ErrCodeInvalidState, err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, app.NewError(app.ErrCodeUnknown, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	if err := s.UserService.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	auth, err := generateAuth(user, s.secret)
	if err != nil {
		return nil, app.NewError(app.ErrCodeUnknown, err)
	}

	if err := saveAuth(ctx, tx, auth); err != nil {
		return nil, mapDBErr(err)
	}
	return auth, nil
}

func generateAuth(user *app.User, secret string) (*app.Auth, error) {
	claims := app.Claims{
		UserID: user.ID,
		Role:   string(user.Role),
	}
	accessToken, err := generateAccessToken(claims, secret)
	if err != nil {
		return nil, app.NewError(app.ErrCodeUnknown, err)
	}
	refreshToken, err := generateRefreshToken(claims, secret)
	if err != nil {
		return nil, app.NewError(app.ErrCodeUnknown, err)
	}

	return &app.Auth{
		UserID:       claims.UserID,
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func saveAuth(ctx context.Context, tx *sql.Tx, auth *app.Auth) error {
	const query = `
	INSERT INTO auth (user_id, access_token, refresh_token)
	VALUES ($1, $2, $3)
	`
	_, err := tx.ExecContext(ctx, query, auth.UserID, auth.AccessToken, auth.RefreshToken)

	return err
}

func updateAuth(ctx context.Context, tx *sql.Tx, auth *app.Auth) error {
	const query = `
	UPDATE auth
	SET access_token = $2, refresh_token = $3
	WHERE user_id = $1
	`
	_, err := tx.ExecContext(ctx, query, auth.UserID, auth.AccessToken, auth.RefreshToken)

	return err
}
func getAuthByUserID(ctx context.Context, tx *sql.Tx, userID string) (*app.Auth, error) {
	const query = `
	SELECT user_id, access_token, refresh_token
	FROM auth
	WHERE user_id = $1
	`

	var auth app.Auth

	row := tx.QueryRowContext(ctx, query, userID)

	err := row.Scan(
		&auth.UserID,
		&auth.AccessToken,
		&auth.RefreshToken,
	)

	if err != nil {
		return nil, err
	}
	return &auth, nil
}

type tokenClaims struct {
	app.Claims
	jwt.RegisteredClaims
}

func generateRefreshToken(claims app.Claims, secret string) (string, error) {
	claims.TokenType = app.RefreshTokenType
	refreshTokenClaims := tokenClaims{
		Claims: claims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(tokenSigningMethod, refreshTokenClaims)

	rawToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return rawToken, nil
}

func generateAccessToken(claims app.Claims, secret string) (string, error) {
	claims.TokenType = app.AccessTokenType
	accessTokenClaims := tokenClaims{
		Claims: claims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(tokenSigningMethod, accessTokenClaims)

	rawToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return rawToken, nil
}
