package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	app "github.com/internships-backend/test-backend-beldurad/internal"
	"github.com/internships-backend/test-backend-beldurad/internal/parser"
)

type ScheduleCreateDTO struct {
	ID         string         `json:"id,omitempty"`
	RoomID     string         `json:"roomId"`
	DaysOfWeek []time.Weekday `json:"daysOfWeek"`
	StartTime  string         `json:"startTime"`
	EndTime    string         `json:"endTime"`
}

func mapCreateDtoToSchedule(log *slog.Logger, dto ScheduleCreateDTO) (*app.Schedule, error) {

	startTime, err := parser.ParseTimeToTodayUTC(dto.StartTime)
	if err != nil {
		return &app.Schedule{}, app.NewError(app.ErrCodeInvalidState, fmt.Errorf("error while parsing:%w", err))
	}
	endTime, err := parser.ParseTimeToTodayUTC(dto.EndTime)
	if err != nil {
		return &app.Schedule{}, app.NewError(app.ErrCodeInvalidState, fmt.Errorf("error while parsing:%w", err))
	}

	result := app.CreateSchedule()
	if dto.ID != "" {
		result.ID = dto.ID
	}
	result.RoomID = dto.RoomID
	result.DaysOfWeek = dto.DaysOfWeek
	result.StartTime = startTime
	result.EndTime = endTime
	return result, nil
}

type ScheduleSaver interface {
	Create(context.Context, *app.Schedule) error
}

type CreateRequest struct {
}

const roomIdParam = "roomId"

func CreateSchedule(log *slog.Logger, scheduleSaver ScheduleSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "schedule.http.schedulehandler"

		roomId := chi.URLParam(r, roomIdParam)

		log = log.With(
			slog.String("op", op),
		)

		var dto ScheduleCreateDTO

		err := render.DecodeJSON(r.Body, &dto)
		if err != nil {
			log.Error("failed to decode request body")
			SendResponseByError(err, w, r)
			return
		}

		if roomId != dto.RoomID {
			log.Error("route param does not equal body param: roomId")
			SendResponseByError(err, w, r)
		}

		schedule, err := mapCreateDtoToSchedule(log, dto)

		if err != nil {
			log.Error("fail during mapping")
			SendResponseByError(err, w, r)
			return
		}

		if err = schedule.Validate(); err != nil {
			log.Error("validation error")
			SendResponseByError(err, w, r)
			return
		}

		err = scheduleSaver.Create(r.Context(), schedule)

		if err != nil {
			log.Error("fail during schedule save")
			SendResponseByError(err, w, r)
			return
		}

		log.Info("new schedule created", slog.Any("schedule", schedule))
		render.Status(r, http.StatusCreated)
	}
}
