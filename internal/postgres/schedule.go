package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/internships-backend/test-backend-beldurad/internal/domain"
)

type ScheduleService struct {
	DB  *sql.DB
	log *slog.Logger
}

func NewScheduleService(db *sql.DB, log *slog.Logger) *ScheduleService {
	return &ScheduleService{
		DB:  db,
		log: log,
	}
}

func (s *ScheduleService) Save(ctx context.Context, schedule *domain.Schedule) error {

	if err := schedule.Validate(); err != nil {
		return err
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", mapDBErr(err))
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	err = saveSchedule(ctx, tx, schedule)

	if err != nil {
		return fmt.Errorf("failed to create schedule: %w", mapDBErr(err))
	}

	err = createNewSlotsBySchedule(ctx, tx, schedule)

	if err != nil {
		return mapDBErr(err)
	}

	return nil
}

func saveSchedule(ctx context.Context, tx *sql.Tx, schedule *domain.Schedule) error {
	const query = `
		INSERT INTO schedule (id, room_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, created_at
	`

	err := tx.QueryRowContext(ctx, query, schedule.ID, schedule.RoomID).
		Scan(&schedule.ID, &schedule.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func getAllSchedules(ctx context.Context, db *sql.DB) ([]*domain.Schedule, error) {
	const query = `
		SELECT id, room_id, created_at, start_time, end_time
		FROM schedule
		ORDER BY created_at DESC
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schedules := make([]*domain.Schedule, 0)

	for rows.Next() {
		var schedule domain.Schedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.RoomID,
			&schedule.CreatedAt,
			&schedule.StartTime,
			&schedule.EndTime,
		)
		if err != nil {
			return nil, err
		}

		weekdays, err := getWeekdaysByScheduleID(ctx, db, schedule.ID)
		if err != nil {
			return nil, err
		}
		schedule.DaysOfWeek = weekdays

		schedules = append(schedules, &schedule)
	}
	return schedules, nil
}

func getWeekdaysByScheduleID(ctx context.Context, db *sql.DB, scheduleID string) ([]time.Weekday, error) {
	const query = `
		SELECT day_of_week
		FROM schedule_day_of_week
		WHERE schedule_id = $1
	`

	rows, err := db.QueryContext(ctx, query, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	weekdays := make([]time.Weekday, 0)
	for rows.Next() {
		var day time.Weekday
		err := rows.Scan(&day)
		if err != nil {
			return nil, err
		}
		weekdays = append(weekdays, day)
	}
	return weekdays, nil
}
