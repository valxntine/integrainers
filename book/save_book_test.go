package book_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/valxntine/integrainers/book"
	"github.com/valxntine/integrainers/entity"
	"testing"
)

type MockGetter struct {
	ErrReturn  error
	BookReturn entity.Book
	Called     int
}

func (bg *MockGetter) GetBook(ctx context.Context, isbn int) (entity.Book, error) {
	bg.Called++
	return bg.BookReturn, bg.ErrReturn
}

type MockSaver struct {
	ErrReturn error
	Called    int
}

func (bs *MockSaver) Save(ctx context.Context, book entity.Book) error {
	bs.Called++
	return bs.ErrReturn
}

func TestSave(t *testing.T) {
	var (
		ctx      context.Context
		testISBN int
	)

	ctx = context.Background()
	testISBN = 123

	t.Run("something fun", func(t *testing.T) {
		mockGetter := &MockGetter{
			ErrReturn: nil,
			BookReturn: entity.Book{
				Author: "Valentine",
				Name:   "Vals Book",
				ISBN:   testISBN,
				Pages:  321,
			},
		}

		mockSaver := &MockSaver{}

		bs := book.New(mockGetter, mockSaver)

		r, err := bs.Save(ctx, testISBN)
		if err != nil {
			t.Fail()
		}

		expectedResponse := book.Response{ISBN: 123}

		assert.Equal(t, expectedResponse, r)
		assert.Equal(t, 1, mockGetter.Called)
		assert.Equal(t, 1, mockSaver.Called)
	})
}
