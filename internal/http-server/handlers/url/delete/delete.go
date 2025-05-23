package deleteUrl

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var alias = chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("empty alias for deleting")
			
			render.JSON(w, r, response.Error("invalid request"))		

			return
		}

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