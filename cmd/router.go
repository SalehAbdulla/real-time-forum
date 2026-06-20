package main

import (
	"net/http"
	"real-time-forum/pkg/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)
	
	mux.Get("/", handlers.Repo.Home)
	return mux
}
