package delete

import (
	"errors"
	"log/slog"
	"net/http"
	resp "rest_api_app/internal/api/responce"
	"rest_api_app/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handelrs.delete.New"

		log = slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := r.URL.Query().Get("alias")
		if alias == "" {
			log.Error("alias is required", sl.Err(errors.New("alias is required")))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("alias is required"))

			return
		}

		err := urlDeleter.DeleteURL(alias)
		if err != nil {
			log.Error("falied to delete url", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("falied to delete url"))

			return
		}

		log.Info("url deleted", slog.String("alias", alias))
		render.JSON(w, r, resp.OK())
	}
}
