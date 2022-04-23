package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/igor6629/booking/internal/config"
	"github.com/igor6629/booking/internal/handlers"
	"github.com/igor6629/booking/internal/models"
	"github.com/igor6629/booking/internal/render"
	"log"
	"net/http"
	"time"
)

const port = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	gob.Register(models.Reservation{})
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Port is %s", port))

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

/*

session.Cookie.Secure = true
app.UseCache = true

*/
