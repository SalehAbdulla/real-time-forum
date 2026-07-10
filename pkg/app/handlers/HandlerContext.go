package handlers

import (
	"net/http"
	"real-time-forum/pkg/config"
	"real-time-forum/pkg/models"
	"real-time-forum/pkg/render"
	"real-time-forum/pkg/app/service"
	pkgwebsocket "real-time-forum/pkg/websocket"
)

var HandlerCtx *HandlerContext

type HandlerContext struct {
	App                *config.AppConfig
	AuthService        service.AuthService
	CategoryService    service.CategoryService
	PostService        service.PostService
	CommentService     service.CommentService
	ReactService       service.ReactionService
	MessageService     service.MessageService
	NotificationService service.NotificationService
	Hub                *pkgwebsocket.Hub
}

func NewHandlerContext(a *config.AppConfig,
	as service.AuthService,
	cs service.CategoryService,
	ps service.PostService,
	cms service.CommentService,
	rs service.ReactionService,
	ms service.MessageService,
	ns service.NotificationService) *HandlerContext {
	return &HandlerContext{
		App:                 a,
		AuthService:         as,
		CategoryService:     cs,
		PostService:         ps,
		CommentService:      cms,
		ReactService:        rs,
		MessageService:      ms,
		NotificationService: ns,
	}
}

func SetHandlerContext(hc *HandlerContext) {
	HandlerCtx = hc
}

func (m *HandlerContext) Home(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderTemplate(w, &models.TemplateData{}); err != nil {
		m.App.Logger.Error("failed to render home template", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
