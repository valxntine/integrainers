package book

import (
	"context"
	"fmt"
	"github.com/valxntine/integrainers/entity"
)

type Getter interface {
	GetBook(ctx context.Context, isbn int) (entity.Book, error)
}

type Saver interface {
	Save(ctx context.Context, book entity.Book) error
}

type Book struct {
	BookGetter Getter
	BookSaver  Saver
}

func New(getter Getter, saver Saver) Book {
	return Book{
		BookGetter: getter,
		BookSaver:  saver,
	}
}

type Response struct {
	ISBN int `json:"isbn"`
}

func (b Book) Save(ctx context.Context, isbn int) (Response, error) {
	bookFromService, err := b.BookGetter.GetBook(ctx, isbn)
	if err != nil {
		return Response{}, fmt.Errorf("getting book: %w", err)
	}

	if err := b.BookSaver.Save(ctx, bookFromService); err != nil {
		return Response{}, fmt.Errorf("saving book: %w", err)
	}
	return Response{ISBN: bookFromService.ISBN}, nil
}
