package main

import (
	"net/http"
	"real-time-forum/pkg/app/handlers"
	pkgmiddleware "real-time-forum/pkg/middleware"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(chiMiddleware.Recoverer)
	mux.Use(RequestLogger)
	mux.Use(SessionLoad)

	mux.Group(func(r chi.Router) {
		r.Post("/api/v1/auth/register", handlers.Repo.Register)
		r.Post("/api/v1/auth/login", handlers.Repo.Login)
	})

	mux.Group(func(r chi.Router) {
		r.Use(pkgmiddleware.AuthMiddleware)

		r.Get("/", handlers.Repo.Home)

		r.Post("/api/v1/auth/logout", handlers.Repo.Logout)
		r.Get("/api/v1/auth/me", handlers.Repo.Me)

		r.Get("/api/v1/categories", handlers.Repo.GetCategories)

		r.Get("/api/v1/posts", handlers.Repo.GetPosts)
		r.Post("/api/v1/posts", handlers.Repo.CreatePost)

		r.Get("/api/v1/posts/comments", handlers.Repo.GetComments)
		r.Post("/api/v1/posts/comments", handlers.Repo.CreateComments)

		r.Post("/api/v1/reactions", handlers.Repo.React)

		r.Get("/api/v1/messages/users", handlers.Repo.GetChatUsers)
		r.Get("/api/v1/messages", handlers.Repo.GetChatMessages)

		r.Get("/ws", handlers.Repo.ServeWs)
	})

	return mux
}
