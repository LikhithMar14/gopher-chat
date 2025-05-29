package api

import (
	"net/http"
	"time"

	"github.com/LikhithMar14/gopher-chat/docs"
	"github.com/LikhithMar14/gopher-chat/internal/config"
	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Application struct {
	Config         config.Config
	Store          store.Storage
	UserService    *service.UserService
	PostService    *service.PostService
	CommentService *service.CommentService
	FollowService  *service.FollowService
	FeedService    *service.FeedService
	AuthService    *service.AuthService
	Version        string
	Logger         *zap.SugaredLogger
}

func NewApplication(cfg config.Config, store store.Storage, version string, logger *zap.SugaredLogger) *Application {
	userService := service.NewUserService(store)
	postService := service.NewPostService(store)
	commentService := service.NewCommentService(store)
	followService := service.NewFollowService(store)
	feedService := service.NewFeedService(store)
	authService := service.NewAuthService(store, cfg.Mail.Exp)

	return &Application{
		Config:         cfg,
		Store:          store,
		UserService:    userService,
		PostService:    postService,
		CommentService: commentService,
		FollowService:  followService,
		FeedService:    feedService,
		AuthService:    authService,
		Version:        version,
		Logger:         logger,
	}
}

func (app *Application) Serve(mux *chi.Mux) error {
	//Docs
	docs.SwaggerInfo.Title = "Gopher Chat API"
	docs.SwaggerInfo.Version = app.Version
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/v1"
	srv := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.Logger.Infow("Server has started", "addr", app.Config.Addr, "env", app.Config.Env, "version", app.Version)

	if err := srv.ListenAndServe(); err != nil {
		app.Logger.Errorw("Failed to start server", "error", err)
		return err
	}

	return nil
}
