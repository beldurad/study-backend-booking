package apphttp

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/internships-backend/test-backend-beldurad/internal/domain"
	"github.com/internships-backend/test-backend-beldurad/internal/logger/sl"
)

const (
	ROOM_CREATE_ENDPOINT = "/rooms/create"
	ROOM_GET_ENDPOINT    = "/rooms/list"
)

type RoomCreateDto struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Capacity    uint   `json:"capacity,omitempty"`
}

func mapCreateDTOToRoom(dto RoomCreateDto) *domain.Room {
	return &domain.Room{
		Name:        dto.Name,
		Description: dto.Description,
		Capacity:    dto.Capacity,
	}
}

type CreateResponse struct {
	domain.Room `json:"room"`
}

type RoomSaver interface {
	Save(ctx context.Context, room *domain.Room) error
}

func SaveRoom(log *slog.Logger, roomSaver RoomSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "room.handler.CreateRoom"

		log = log.With(
			slog.String("op", op),
		)

		var dto RoomCreateDto

		err := render.DecodeJSON(r.Body, &dto)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			SendResponseByError(err, w, r)
			return
		}

		roomToSave := mapCreateDTOToRoom(dto)
		if err := roomToSave.Validate(); err != nil {
			SendResponseByError(err, w, r)
		}

		err = roomSaver.Save(r.Context(), roomToSave)
		if err != nil {
			log.Error("fail during room save", sl.Err(err))
			SendResponseByError(err, w, r)
			return
		}
		response := CreateResponse{Room: *roomToSave}

		log.Info("new room created", slog.Any("room", *roomToSave))
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, response)
	}
}

type RoomGetter interface {
	GetAll(ctx context.Context) ([]*domain.Room, error)
}

type GetResponse struct {
	Rooms []*domain.Room `json:"rooms"`
}

func GetRooms(log *slog.Logger, roomGetter RoomGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handler.room.GetRooms"

		log = log.With(
			slog.String("op", op),
		)

		rooms, err := roomGetter.GetAll(r.Context())
		if err != nil {
			log.Error("failed to get rooms", sl.Err(err))
			SendResponseByError(err, w, r)
			return
		}

		if len(rooms) == 0 {
			log.Info("no rooms found")
			render.Status(r, http.StatusNoContent)
			return
		}

		response := GetResponse{
			Rooms: make([]*domain.Room, 0, len(rooms)),
		}

		for _, room := range rooms {
			response.Rooms = append(response.Rooms, room)
		}

		log.Info("rooms retrieved successfully", slog.Int("count", len(response.Rooms)))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)
	}
}
