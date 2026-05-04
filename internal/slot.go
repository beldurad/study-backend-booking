package app

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	SlotStatusAvailable   = "available"
	SlotStatusUnavailable = "unavailable"
)

type SlotStatus string

func (s SlotStatus) Validate() error {
	if s != SlotStatusAvailable && s != SlotStatusUnavailable {
		return NewError(ErrCodeInvalidState, fmt.Errorf("invalid slot status"))
	}
	return nil
}

type Slot struct {
	ID        string     `db:"id"`
	RoomID    string     `db:"room_id"`
	StartTime time.Time  `db:"start_time"`
	EndTime   time.Time  `db:"end_time"`
	Status    SlotStatus `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
}

func CreateSlot() *Slot {
	s := &Slot{}
	s.ID = uuid.NewString()
	s.CreatedAt = time.Now()
	return s
}

func (s *Slot) Validate() error {
	if uuid.Validate(s.RoomID) != nil {
		return NewError(ErrCodeInvalidState, fmt.Errorf("room id need to be valid uuid"))
	}
	if uuid.Validate(s.ID) != nil {
		return NewError(ErrCodeInvalidState, fmt.Errorf("id is required and need to be valid uuid"))
	}
	if s.Status.Validate() != nil {
		return NewError(ErrCodeInvalidState, fmt.Errorf("invalid slot status"))
	}
	if s.CreatedAt.IsZero() {
		return NewError(ErrCodeInvalidState, fmt.Errorf("createdAt is required"))
	}
	return nil
}
