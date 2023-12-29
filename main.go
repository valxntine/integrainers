package integrainers

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/valxntine/integrainers/app"
	"github.com/valxntine/integrainers/book"
	"github.com/valxntine/integrainers/config"
	"github.com/valxntine/integrainers/services/library"
	"log"
	"net/http"
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
	httpClient := http.DefaultClient

	libraryClient := library.New(*httpClient, cfg.LibraryHost)

	book := book.New(libraryClient, repo)

	application := app.App{
		Config:    cfg,
		Router:    mux.NewRouter(),
		BookSaver: book,
	}

	if err := application.Routes(); err != nil {
		log.Fatal(err)
	}
	application.Run()
}
