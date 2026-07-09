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
		r.Post("/api/v1/auth/register", handlers.HandlerCtx.Register)
		r.Post("/api/v1/auth/login", handlers.HandlerCtx.Login)
	})

	mux.Group(func(r chi.Router) {
		r.Use(pkgmiddleware.AuthMiddleware)

		r.Get("/", handlers.HandlerCtx.Home)

		r.Post("/api/v1/auth/logout", handlers.HandlerCtx.Logout)
		r.Get("/api/v1/auth/me", handlers.HandlerCtx.Me)

		r.Get("/api/v1/categories", handlers.HandlerCtx.GetCategories)

		r.Get("/api/v1/posts", handlers.HandlerCtx.GetPosts)
		r.Post("/api/v1/posts", handlers.HandlerCtx.CreatePost)

		r.Get("/api/v1/posts/comments", handlers.HandlerCtx.GetComments)
		r.Post("/api/v1/posts/comments", handlers.HandlerCtx.CreateComments)

		r.Post("/api/v1/reactions", handlers.HandlerCtx.React)

		r.Get("/api/v1/messages/users", handlers.HandlerCtx.GetChatUsers)
		r.Get("/api/v1/messages", handlers.HandlerCtx.GetChatMessages)

		r.Get("/ws", handlers.HandlerCtx.ServeWs)
	})

	return mux
}
