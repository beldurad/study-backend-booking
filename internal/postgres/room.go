package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	app "github.com/internships-backend/test-backend-beldurad/internal"
)

type RoomService struct {
	DB  *sql.DB
	log *slog.Logger
}

func NewRoomService(db *sql.DB, log *slog.Logger) *RoomService {
	return &RoomService{
		DB:  db,
		log: log,
	}
}

func (s *RoomService) Save(ctx context.Context, room *app.Room) error {
	const op = "room.repository.CreateRoom"
	const query = `
        INSERT INTO room (name, description, capacity, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id, name, description, capacity, created_at
    `

	if err := room.Validate(); err != nil {
		return fmt.Errorf("%v: failed to validate room: %w", op, err)
	}

	err := s.DB.QueryRowContext(
		ctx,
		query,
		room.Name,
		room.Description,
		room.Capacity,
		room.CreatedAt,
	).Scan(
		&room.ID,
		&room.Name,
		&room.Description,
		&room.Capacity,
		&room.CreatedAt,
	)

	if err != nil {
		err = mapDBErr(err)
		return fmt.Errorf("%v: failed to save room in db: %w", op, err)
	}

	return nil
}

func (s *RoomService) GetAll(ctx context.Context) ([]*app.Room, error) {
	const op = "room.repository.GetAllRooms"

	query := `
		SELECT id, name, description, capacity, created_at
		FROM rooms
		ORDER BY created_at DESC
	`

	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		err = mapDBErr(err)
		return nil, fmt.Errorf("%v: select query failed for rooms failed: %w", op, err)
	}
	defer rows.Close()

	var rooms []*app.Room
	for rows.Next() {
		var room app.Room
		err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.Description,
			&room.Capacity,
			&room.CreatedAt,
		)
		if err != nil {
			return nil, mapDBErr(err)
		}
		rooms = append(rooms, &room)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return rooms, nil
}
