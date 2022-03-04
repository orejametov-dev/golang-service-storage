package router

import (
	"net/http"
	"orejametov/service-storage/internal/config"
	"orejametov/service-storage/internal/file"
	"orejametov/service-storage/pkg/logging"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(cfg *config.Config, fileService file.Service) http.Handler {
	var (
		logger       = logging.GetLogger()
		logFormatter = middleware.DefaultLogFormatter{Logger: logger}
	)

	var (
		logging       = middleware.RequestLogger(&logFormatter)
	)

	router := chi.NewRouter()
	router.Use(CORS)

	router.With(logging, middleware.Recoverer).Route("/api/files", func(r chi.Router) {
		r.Get("/", file.GetFilesByNoteUUID(fileService, logger))
		r.Get("/{id}", file.GetFile(fileService, logger))
		r.Post("/", file.CreateFile(fileService, logger))
		r.Get("/{id}", file.DeleteFile(fileService, logger))
	})

	return router
}
