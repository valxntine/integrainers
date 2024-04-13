package book

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/valxntine/integrainers/entity"
	"strconv"
)

//var ErrorDatabaseConnectionFailed = errors.New("database connection failed")
//var ErrorInsertFailed = errors.New("failed to insert record")
//var ErrorUpdateFailed = errors.New("failed to update record")
//var ErrorInvalidEanArgument = errors.New("the ean supplied isn't valid")
//var ErrorPrepareStmtFailed = errors.New("preparing sql query has failed")
//
//// TODO: Make sure that ErrorItemNotFound is returned by some code path here in the repository, and then handled
//// in the business layer and the handler.
//var ErrorItemNotFound = errors.New("the item wasn't found")
//var ErrorQueryFailure = errors.New("failed while querying item")
//var ErrorInsertItem = errors.New("failed inserting item")
//var ErrorNoRows = errors.New("rescan was not found")

type Repository struct {
	db          *sql.DB
	credentials string
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) Save(ctx context.Context, book entity.Book) error {
	query := `
		INSERT INTO book (isbn, author, name, pages) values (?, ?, ?, ?)`

	isbnStr := strconv.Itoa(book.ISBN)

	_, err := r.db.ExecContext(ctx, query, isbnStr, book.Author, book.Name, book.Pages)
	if err != nil {
		return fmt.Errorf("exec context: %w", err)
	}
	return nil
}
