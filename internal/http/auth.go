package apphttp

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/internships-backend/test-backend-beldurad/internal/domain"
	"github.com/internships-backend/test-backend-beldurad/internal/logger/sl"
)

type dummyLoginRequest struct {
	Role string `json:"role"`
}

type UserGetter interface {
	GetDummyUser(role string) (*domain.User, error)
}

type TokenGenerator interface {
	Generate(userID string, role string) (string, error)
}

type Response struct {
	Token string `json:"token"`
}

func DummyLogin(log *slog.Logger, userGetter UserGetter, tokenGen TokenGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dummyLoginRequest

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed", sl.Err(err))

			return
		}

		dummyUser, err := userGetter.GetDummyUser(req.Role)
		if err != nil {
			SendResponseByError(err, w, r)
			return
		}

		token, err := tokenGen.Generate(dummyUser.Id, dummyUser.Role)
		if err != nil {
			log.Error("token generation error", sl.Err(err))
			SendResponseByError(err, w, r)
			return
		}
		log.Debug("token", slog.String("jwt", token))

		render.JSON(w, r, Response{Token: token})
	}
}
