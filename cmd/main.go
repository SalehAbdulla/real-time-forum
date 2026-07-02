package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"real-time-forum/pkg/config"
	"real-time-forum/pkg/handlers"
	"real-time-forum/pkg/render"
	db "real-time-forum/pkg/repos"
	"real-time-forum/pkg/service"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":3000"
var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false
	app.UseCache = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	app.TemplateCache = templateCache

	render.NewTemplates(&app)

	// Initialize database
	database, err := sql.Open("sqlite3", "./pkg/repos/realTimeForum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	dbConn := &db.DB{Conn: database}

	// Initialize services
	authService := service.NewAuthService(dbConn)

	repo := handlers.NewRepo(&app, authService)
	handlers.NewHandlers(repo)

	fmt.Printf("Starting application on port %s\n", portNumber)

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}

	err = serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}