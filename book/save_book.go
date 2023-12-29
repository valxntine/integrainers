package book

import (
	"context"
	"github.com/valxntine/integrainers/entity"
)

type Getter interface {
	GetBook(ctx context.Context, iban string) (entity.Book, error)
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
	Iban string `json:"iban"`
}

func (b Book) Save(ctx context.Context, iban string) (Response, error) {
	bookFromService, err := b.BookGetter.GetBook(ctx, iban)
	if err != nil {
		return Response{}, err
	}

	if err := b.BookSaver.Save(ctx, bookFromService); err != nil {
		return Response{}, err
	}
	return Response{Iban: bookFromService.Iban}, nil
}
