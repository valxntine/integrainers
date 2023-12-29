package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/valxntine/integrainers/config"
	"github.com/valxntine/integrainers/handlers"
	"log"
	"net/http"
	"time"
)

type App struct {
	Config    config.AppConfig
	Router    *mux.Router
	BookSaver handlers.Saver
}

func (a *App) Run() {
	const forDockerUseGlobalIP = "0.0.0.0"
	fmt.Printf("App is running on port: %s\n", a.Config.Port)

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%s", forDockerUseGlobalIP, a.Config.Port),
		Handler:      a.Router,
		IdleTimeout:  65 * time.Second,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func (a *App) Routes() error {
	a.Router.Handle("/api/v1/books/", handlers.SaveBook(a.BookSaver)).Methods(http.MethodPost)
	return nil
}
