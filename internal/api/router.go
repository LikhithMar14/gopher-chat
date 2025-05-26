package api

import (
	"time"

	"github.com/LikhithMar14/gopher-chat/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	healthHandler := handlers.NewHealthHandler(app.Config)
	userHandler := handlers.NewUserHandler(app.UserService)

	r.Get("/v1/health", healthHandler.Handle)

	r.Route("/v1/users", func(r chi.Router) {
		r.Get("/", userHandler.GetUsers)
		r.Post("/", userHandler.CreateUser)
	})

	return r
}
