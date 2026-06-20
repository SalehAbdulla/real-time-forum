package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"real-time-forum/pkg/config"
	"real-time-forum/pkg/handlers"
	"real-time-forum/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":3000"
var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false
	app.UseCache = false;

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {log.Fatal(err)}
	app.TemplateCache = templateCache

	render.NewTemplates(&app)

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)


	fmt.Printf("Starting applicaiton on port %s\n", portNumber)

	serve := &http.Server{
		Addr: portNumber,
		Handler: routes(),
	}

	err = serve.ListenAndServe()
}
