package main

import (
	"log"
	"net/http"
	"time"

	"github.com/LikhithMar14/gopher-chat/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store store.Storage

}

type config struct {
	addr string
	db dbConfig

}

type dbConfig struct {
	addr string
	maxOpenConns int
	maxIdleConns int
	maxLifetime string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/v1/health", app.healthcheckHandler)

	return r
}

func (app *application) run(mux *chi.Mux) error {
	srv := &http.Server{
		Addr: app.config.addr,
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout: time.Minute,
	}

	log.Printf("Starting server on %s", app.config.addr)

	return srv.ListenAndServe()
}