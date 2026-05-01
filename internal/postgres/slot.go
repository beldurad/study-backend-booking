package postgres

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/internships-backend/test-backend-beldurad/internal/domain"
)

const timeBound = 10 * 24 * time.Hour

type SlotWorker struct {
	db  *sql.DB
	log *slog.Logger
}

func NewSlotService(db *sql.DB, log *slog.Logger) *SlotWorker {
	return &SlotWorker{
		db:  db,
		log: log,
	}
}

func (s *SlotWorker) Start(ctx context.Context) error {
	ticker := time.NewTicker(6 * time.Hour)

	for {
		select {
		case <-ticker.C:
			s.CreateSlotsForAllSchedules(ctx)
		case <-ctx.Done():
			ticker.Stop()
			return nil
		}
	}
}

func (s *SlotWorker) CreateSlotsForAllSchedules(ctx context.Context) {

	schedules, err := getAllSchedules(ctx, s.db)
	if err != nil {
		return
	}

	for _, schedule := range schedules {
		func() {
			tx, err := s.db.BeginTx(ctx, nil)
			if err != nil {
				return
			}
			defer func() {
				if err != nil {
					tx.Rollback()
				}
				tx.Commit()
			}()
			createNewSlotsBySchedule(ctx, tx, schedule)
		}()
	}
}

func createNewSlotsBySchedule(ctx context.Context, tx *sql.Tx, schedule *domain.Schedule) error {

	lastDate, err := getLastSlotDate(ctx, tx, schedule.RoomID)
	if err != nil {
		if err == sql.ErrNoRows {
			lastDate = time.Now().Add(10 * time.Minute)
		} else {
			return err
		}
	}

	for cur := lastDate; cur.Before(lastDate.Add(timeBound)); {

		err = createSlotsForDay(ctx, tx, schedule, cur)

		if err != nil {
			return err
		}

		cur = cur.Add(24 * time.Hour).Truncate(24 * time.Hour)
	}

	return nil

}

func createSlotsForDay(ctx context.Context, tx *sql.Tx, schedule *domain.Schedule, day time.Time) error {
	if !schedule.ContainsWeekday(day.Weekday()) {
		return nil
	}

	startTime := time.Date(
		day.Year(),
		day.Month(),
		day.Day(),
		schedule.StartTime.Hour(),
		schedule.StartTime.Minute(),
		0, 0, time.UTC,
	)

	endTime := time.Date(
		day.Year(),
		day.Month(),
		day.Day(),
		schedule.EndTime.Hour(),
		schedule.EndTime.Minute(),
		0, 0, time.UTC,
	)

	if endTime.Before(startTime) {
		endTime = endTime.Add(24 * time.Hour)
	}

	if endTime.Before(day) {
		return nil
	}

	cur := startTime

	if cur.Before(day) {
		cur = day.Truncate(time.Second)
	}

	slots := make([]*domain.Slot, 10)

	for ; cur.Before(endTime); cur = cur.Add(30 * time.Minute) {
		slotToAdd := domain.CreateSlot()
		slotToAdd.RoomID = schedule.RoomID
		slotToAdd.StartTime = cur
		slotToAdd.EndTime = cur.Add(30 * time.Minute)
		slotToAdd.Status = domain.SlotStatusAvailable
		slots = append(slots, slotToAdd)
	}

	return saveAllSlots(ctx, tx, slots)

}

func getLastSlotDate(ctx context.Context, tx *sql.Tx, roomID string) (time.Time, error) {
	const query = `
		SELECT end_time
		FROM slot
		WHERE room_id = $1
		ORDER BY end_time DESC
		LIMIT 1
	`
	var result time.Time

	if err := tx.QueryRowContext(ctx, query, roomID).Scan(&result); err != nil {
		return time.Time{}, err
	}

	return result, nil
}

func saveAllSlots(ctx context.Context, tx *sql.Tx, slots []*domain.Slot) error {
	const query = `
		INSERT INTO slot (id, room_id, start_time, end_time, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	for _, slot := range slots {
		if _, err := tx.ExecContext(ctx, query, slot.ID, slot.RoomID, slot.StartTime, slot.EndTime, slot.Status, slot.CreatedAt); err != nil {
			return err
		}
	}
	return nil
}
