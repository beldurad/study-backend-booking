package postgres

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	app "github.com/internships-backend/test-backend-beldurad/internal"
)

type UserService struct {
	app.PasswordHasher
	DB  *sql.DB
	log *slog.Logger
}

func NewUserService(db *sql.DB, hasher app.PasswordHasher, log *slog.Logger) *UserService {
	return &UserService{
		PasswordHasher: hasher,
		DB:             db,
		log:            log,
	}
}

func (s *UserService) GetDummyUser(ctx context.Context, role app.Role) (*app.User, error) {
	if err := role.Validate(); err != nil {
		return nil, app.NewError(app.ErrCodeInvalidState, err)
	}
	var res *app.User
	switch role {
	case app.AdminRole:
		res = &app.User{
			ID:        "1",
			Email:     "admin@test",
			Role:      app.AdminRole,
			CreatedAt: time.Now(),
		}
		return res, nil
	case "user":
		res = &app.User{
			ID:        "2",
			Email:     "user@test",
			Role:      app.UserRole,
			CreatedAt: time.Now(),
		}
		return res, nil
	default:
		return nil, app.NewError(app.ErrCodeInvalidState, nil)
	}
}

func (s *UserService) GetUser(ctx context.Context, id string) (*app.User, error) {
	if err := uuid.Validate(id); err != nil {
		return nil, app.NewError(app.ErrCodeInvalidState, err)
	}
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, app.NewError(app.ErrCodeUnknown, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	user, err := getUserById(ctx, tx, id)

	if err != nil {
		return nil, mapDBErr(err)
	}
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *app.User) error {
	if user.PasswordHash == "" {
		hash, err := s.PasswordHasher.Hash(user.Password)
		if err != nil {
			return app.NewError(app.ErrCodeUnknown, err)
		}
		user.PasswordHash = hash
	}
	if err := user.Validate(); err != nil {
		return app.NewError(app.ErrCodeInvalidState, err)
	}
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return app.NewError(app.ErrCodeUnknown, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()
	if err := createUser(ctx, tx, user); err != nil {
		return mapDBErr(err)
	}
	return nil
}

func createUser(ctx context.Context, tx *sql.Tx, user *app.User) error {
	const query = `
	INSERT INTO users (id, email, role, created_at)
	VALUES ($1, $2, $3, $4)
	`
	_, err := tx.ExecContext(ctx, query, user.ID, user.Email, user.Role, user.CreatedAt)
	return err
}
func getUserById(ctx context.Context, tx *sql.Tx, id string) (*app.User, error) {
	const query = `
		SELECT id, email, role, created_at
		FROM users
		WHERE id = $1
	`
	row := tx.QueryRowContext(ctx, query, id)
	var u app.User
	err := row.Scan(&u.ID, &u.Email, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
