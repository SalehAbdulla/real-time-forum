package handlers

import (
	"log"
	"net/http"
	"real-time-forum/pkg/config"
	"real-time-forum/pkg/models"
	"real-time-forum/pkg/render"
	"real-time-forum/pkg/service"
)

var Repo *Repository

type Repository struct {
	App         *config.AppConfig
	AuthService service.AuthService
}

func NewRepo(a *config.AppConfig, as service.AuthService) *Repository {
	return &Repository{
		App:         a,
		AuthService: as,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	if err := render.RenderTemplate(w, &models.TemplateData{}); err != nil {
		log.Fatal(err.Error())
	}
}