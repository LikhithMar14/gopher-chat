package api

import (
	"log"
	"net/http"
	"time"

	"github.com/LikhithMar14/gopher-chat/internal/config"
	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/store"
	"github.com/go-chi/chi/v5"
)

type Application struct {
	Config      config.Config
	Store       store.Storage
	UserService *service.UserService
	PostService *service.PostService
	CommentService *service.CommentService
}

func NewApplication(cfg config.Config, store store.Storage) *Application {
	userService := service.NewUserService(store)
	postService := service.NewPostService(store)
	commentService := service.NewCommentService(store)
	return &Application{
		Config:      cfg,
		Store:       store,
		UserService: userService,
		PostService: postService,
		CommentService: commentService,
	}
}

func (app *Application) Serve(mux *chi.Mux) error {
	srv := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,	
	}

	log.Printf("Starting server on %s", app.Config.Addr)

	return srv.ListenAndServe()
}
