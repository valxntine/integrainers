package handlers

import (
	"context"
	"encoding/json"
	"github.com/valxntine/integrainers/book"
	"net/http"
)

type Saver interface {
	Save(ctx context.Context, iban string) (book.Response, error)
}

type Request struct {
	Iban string `json:"iban"`
}

type Response struct {
	Message string `json:"message"`
	Iban    string `json:"iban"`
}

func SaveBook(saver Saver) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var body Request
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid iban"`))
			return
		}

		res, err := saver.Save(ctx, body.Iban)
		if err != nil {
			w.WriteHeader(http.StatusTeapot)
		}
		w.WriteHeader(http.StatusCreated)
		bytes, _ := json.Marshal(Response{
			Message: "Successfully saved book",
			Iban:    res.Iban,
		})
		w.Write(bytes)
	})
}
