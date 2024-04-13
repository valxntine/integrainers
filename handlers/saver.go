package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/valxntine/integrainers/book"
	"net/http"
)

type Saver interface {
	Save(ctx context.Context, isbn int) (book.Response, error)
}

type Request struct {
	ISBN int `json:"isbn"`
}

type Response struct {
	Message string `json:"message"`
	ISBN    int    `json:"isbn"`
}

func SaveBook(saver Saver) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var body Request
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid isbn"}`))
			return
		}

		res, err := saver.Save(ctx, body.ISBN)
		if err != nil {
			w.WriteHeader(http.StatusTeapot)
			w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, err)))
			return
		}
		w.WriteHeader(http.StatusCreated)
		bytes, _ := json.Marshal(Response{
			Message: "Successfully saved book",
			ISBN:    res.ISBN,
		})
		w.Write(bytes)
	})
}
