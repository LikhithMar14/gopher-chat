package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
	"github.com/LikhithMar14/gopher-chat/internal/handlers"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)




func (app *Application) postsContextMiddleware(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.ReadIDParam(r)
		if err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		ctx := r.Context()
		post, err := app.PostService.GetPostByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err,apperrors.ErrPostNotFound):
					utils.WriteJSONError(w, http.StatusNotFound, err.Error())
					return
				default:
					utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
					return	
			}
		}
		ctx = context.WithValue(ctx, utils.PostIDKey, post)
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
	r.Route("/v1", func(r chi.Router) {

		r.Get("/health", healthHandler.Handle)

		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.GetUsers)
			r.Post("/", userHandler.CreateUser)
		})

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", postHandler.CreatePost)
			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)		
				r.Get("/", postHandler.GetPostByID)
				r.Delete("/", postHandler.DeletePost)
				r.Patch("/", postHandler.UpdatePost)
			})
		})
	})


	return r
}
