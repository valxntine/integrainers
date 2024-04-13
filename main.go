package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/valxntine/integrainers/app"
	"github.com/valxntine/integrainers/book"
	"github.com/valxntine/integrainers/config"
	"github.com/valxntine/integrainers/services/library"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.NewAppConfig()
	if err != nil {
		log.Fatal(err)
	}
	conn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatal(err)
	}

	repo := book.NewRepository(db)
	httpClient := http.Client{
		Transport: http.DefaultTransport,
		Timeout:   30 * time.Second,
	}

	libraryClient := library.New(httpClient, cfg.LibraryHost)

	bookSvc := book.New(libraryClient, repo)

	application := app.App{
		Config:    cfg,
		Router:    mux.NewRouter(),
		BookSaver: bookSvc,
	}

	if err = application.Routes(); err != nil {
		log.Fatal(err)
	}
	application.Run()
}
