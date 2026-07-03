package main

import (
	"net/http"
	"real-time-forum/pkg/handlers"
	pkgmiddleware "real-time-forum/pkg/middleware"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(chiMiddleware.Recoverer)
	mux.Use(RequestLogger)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Post("/api/v1/auth/register", handlers.Repo.Register)
	mux.Post("/api/v1/auth/login", handlers.Repo.Login)
	mux.With(pkgmiddleware.AuthMiddleware).Post("/api/v1/auth/logout", handlers.Repo.Logout)
	mux.With(pkgmiddleware.AuthMiddleware).Get("/api/v1/auth/me", handlers.Repo.Me)
	return mux
}
