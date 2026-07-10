package main

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"real-time-forum/pkg/config"
	"real-time-forum/pkg/app/handlers"
	"real-time-forum/pkg/logger"
	"real-time-forum/pkg/render"
	db "real-time-forum/pkg/app/repositories"
	"real-time-forum/pkg/app/service"
	pkgwebsocket "real-time-forum/pkg/websocket"
)

const portNumber = ":3000"

var app config.AppConfig

func main() {
	app.InProduction = false
	app.UseCache = false
	app.LogLevel = os.Getenv("LOG_LEVEL")

	// Initialize structured logger
	logger.InitLogger(&app)

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		app.Logger.Error("failed to create template cache", "error", err)
		os.Exit(1)
	}
	app.TemplateCache = templateCache

	render.NewTemplates(&app)

	// Initialize database
	database, err := sql.Open("sqlite3", "./pkg/app/repositories/realTimeForum.db")
	if err != nil {
		app.Logger.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	dbConn := &db.DB{Conn: database}

	// Initialize services
	authService := service.NewAuthService(dbConn)
	categoryService := service.NewCategoryService(dbConn)
	postService := service.NewPostService(dbConn)
	commentService := service.NewCommentService(dbConn)
	reactService := service.NewReactionService(dbConn)
	messageService := service.NewMessageService(dbConn, dbConn)

	hc := handlers.NewHandlerContext(&app, authService, categoryService, postService, commentService, reactService, messageService)
	handlers.SetHandlerContext(hc)

	// Initialize WebSocket hub
	wsHub := pkgwebsocket.NewHub()
	hc.SetHub(wsHub)
	go wsHub.Run()

	app.Logger.Info("starting application", "port", portNumber)

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}

	err = serve.ListenAndServe()
	if err != nil {
		app.Logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}