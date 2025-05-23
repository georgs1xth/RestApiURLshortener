package deleteUrl

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/utils"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) error
	GetURL(alias string) (string, error)
}

type Request struct {
	Alias string `json:"alias" validate:"required"`
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Info("failed to decode request body",)

			render.JSON(w, r, response.Error(fmt.Sprintf("failed to decode request: %v", err)))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := utils.Validate.Struct(req); err != nil {
			log.Error("invalid alias field")

			render.JSON(w, r, response.Error("alias field is required"))

			return
		}

		alias := req.Alias

		if err := urlDeleter.DeleteURL(alias); err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("the url doesn't exist", sl.Err(err))
				
				render.JSON(w, r, response.Error("the url doesn't exist"))
				
				return 
			}
			log.Error("error deleting url", sl.Err(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		log.Info(fmt.Sprintf("successfully deleted url for alias: %s", alias))

		render.JSON(w, r, response.OK())
	}
}