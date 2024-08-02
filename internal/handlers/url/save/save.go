package save

import (
	"log/slog"
	"net/http"
	resp "rest_api_app/internal/api/responce"
	"rest_api_app/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"` // для пакета validator
	Alias string `json:"alias,omitempty"`             // omitempty это парамент в json он говорит о том что если этот параметр пустой то в итогом
}

type Response struct {
	resp.Responce
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "hadnlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("falied to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("falied to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}
	}
}
