package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "rest_api_app/internal/api/responce"
	"rest_api_app/internal/lib/logger/sl"
	"rest_api_app/internal/lib/random"
	"rest_api_app/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"` // для пакета validator
	Alias string `json:"alias,omitempty"`             // omitempty это парамент в json он говорит о том что если этот параметр пустой то в итогом
}

// TODO move to config
const aliasLength = 6

type Response struct {
	resp.Responce
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
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

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("falied to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		// Обработка ошибки если такой алиас уже есть
		if errors.Is(err, storage.ErrUrlExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.Status(r, http.StatusConflict)
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}
		// просто обработка ошибки
		if err != nil {
			log.Error("falied to add url", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("falied to add url"))
			return
		}
		log.Info("url added", slog.Int64("id", id))

		responceOk(w, r, alias)
	}
}

func responceOk(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Responce: resp.OK(),
		Alias:    alias,
	})
}
