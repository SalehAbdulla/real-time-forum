package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"real-time-forum/pkg/config"
	"real-time-forum/pkg/handlers"
	"real-time-forum/pkg/logger"
	"real-time-forum/pkg/render"
	db "real-time-forum/pkg/repositories"
	"real-time-forum/pkg/service"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":3000"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false
	app.UseCache = false
	app.LogLevel = os.Getenv("LOG_LEVEL")

	// Initialize structured logger
	logger.InitLogger(&app)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		app.Logger.Error("failed to create template cache", "error", err)
		os.Exit(1)
	}
	app.TemplateCache = templateCache

	render.NewTemplates(&app)

	// Initialize database
	database, err := sql.Open("sqlite3", "./pkg/repositories/realTimeForum.db")
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
	messageService := service.NewMessageService(dbConn)

	repo := handlers.NewRepo(&app, authService, categoryService, postService, commentService, reactService, messageService)
	handlers.NewHandlers(repo)

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
