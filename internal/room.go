package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name" validate:"required"`
	Description string    `db:"description" json:"description"`
	Capacity    uint      `db:"capacity" json:"capacity"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

func CreateRoom() *Room {
	return &Room{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
	}
}

func (r *Room) Validate() error {
	if r.Name == "" {
		return NewError(ErrCodeInvalidState, fmt.Errorf("name is required"))
	}
	if uuid.Validate(r.ID) != nil {
		return NewError(ErrCodeInvalidState, fmt.Errorf("id is required and need to be valid uuid"))
	}
	if r.CreatedAt.IsZero() {
		return NewError(ErrCodeInvalidState, fmt.Errorf("createdAt is required"))
	}
	return nil
}

type Service interface {
	Save(context.Context, *Room) error
	GetAll(context.Context) ([]*Room, error)
}
