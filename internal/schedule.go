package app

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID         string         `json:"id" db:"id"`
	RoomID     string         `json:"roomId" db:"room_id"`
	CreatedAt  time.Time      `json:"createdAt" db:"created_at"`
	StartTime  time.Time      `json:"startTime" db:"start_time"`
	EndTime    time.Time      `json:"endTime" db:"end_time"`
	DaysOfWeek []time.Weekday `json:"daysOfWeek" db:"-"`
}

func CreateSchedule() *Schedule {
	res := &Schedule{}
	res.ID = uuid.NewString()
	res.CreatedAt = time.Now()
	return res
}

func (s *Schedule) Validate() error {
	if uuid.Validate(s.RoomID) != nil {
		return NewError(ErrCodeInvalidState, fmt.Errorf("room id need to be valid uuid"))
	}
	if uuid.Validate(s.ID) != nil {
		return NewError(ErrCodeInvalidState, fmt.Errorf("id is required and need to be valid uuid"))
	}
	if s.CreatedAt.IsZero() {
		return NewError(ErrCodeInvalidState, fmt.Errorf("createdAt is required"))
	}
	if s.StartTime.IsZero() {
		return NewError(ErrCodeInvalidState, fmt.Errorf("startTime is required"))
	}
	if s.EndTime.IsZero() {
		return NewError(ErrCodeInvalidState, fmt.Errorf("endTime is required"))
	}
	return nil
}

func (s *Schedule) ContainsWeekday(day time.Weekday) bool {
	return slices.Contains(s.DaysOfWeek, day)
}

type ScheduleService interface {
	Save(ctx context.Context, schedule *Schedule) error
}
