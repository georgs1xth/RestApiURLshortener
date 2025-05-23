package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.4 --name URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		// Init lines for logger
		const op = "handlers.redirect.New"
		
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		url, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {			
				log.Info("url not found", "alias", alias)
	
				render.JSON(w, r, response.Error("not found"))

				return
			}

			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, response.Error("internal error"))
			
			return
		}

		log.Info("got url", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusFound)
	}	
}