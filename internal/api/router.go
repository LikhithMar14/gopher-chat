package api

import (
	"context"
	"errors"

	"net/http"
	"time"

	"github.com/LikhithMar14/gopher-chat/internal/handlers"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.ReadIDParam(r)

		if err != nil {
			utils.HandleValidationError(w, errors.New("invalid input format"))
			return
		}
		ctx := r.Context()
		post, err := app.PostService.GetPostByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, apperrors.ErrPostNotFound):
				utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			default:
				utils.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}
		ctx = context.WithValue(ctx, utils.PostIDKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (app *Application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := utils.ReadIDParam(r)
		if err != nil {
			utils.HandleValidationError(w, errors.New("invalid input format"))
			return
		}

		user, err := app.UserService.GetUserByID(ctx, id)

		if err != nil {
			//will give internal server error if the userid is not correct
			utils.HandleInternalError(w, err)
			return
		}

		ctx = context.WithValue(ctx, utils.UserIDKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (app *Application) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	healthHandler := handlers.NewHealthHandler(app.Config)
	userHandler := handlers.NewUserHandler(app.UserService)
	postHandler := handlers.NewPostHandler(app.PostService, app.CommentService)
	commentHandler := handlers.NewCommentHandler(app.CommentService, app.PostService)
	followHandler := handlers.NewFollowHandler(app.FollowService, app.UserService)
	feedHandler := handlers.NewFeedHandler(app.UserService, app.PostService, app.FeedService)

	r.Route("/v1", func(r chi.Router) {

		r.Get("/health", healthHandler.Handle)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", postHandler.CreatePost)
			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", postHandler.GetPostByID)
				r.Delete("/", postHandler.DeletePost)
				r.Patch("/", postHandler.UpdatePost)
				r.Route("/comments", func(r chi.Router) {
					r.Post("/", commentHandler.CreateComment)
					r.Get("/", commentHandler.GetCommentsByPostID)
				})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.GetUsers)
			r.Post("/", userHandler.CreateUser)
			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Get("/", userHandler.GetUserByID)
				r.Put("/follow", followHandler.FollowUser)
				r.Put("/unfollow", followHandler.UnfollowUser)
			})
			r.Route("/feed", func(r chi.Router) {
					r.Get("/", feedHandler.GetFeed)
			})

		})
	})

	return r
}
