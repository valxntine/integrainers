package library

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valxntine/integrainers/entity"
	"net/http"
)

var (
	ErrBookNotFound       = errors.New("book not found")
	ErrUnexpectedResponse = errors.New("unexpected response")
	ErrInvalidResponse    = errors.New("invalid response")
)

type Client struct {
	HTTPClient http.Client
	BaseURL    string
}

func New(client http.Client, url string) Client {
	return Client{
		HTTPClient: client,
		BaseURL:    url,
	}
}

func (c Client) GetBook(ctx context.Context, isbn int) (entity.Book, error) {
	b := entity.Book{}

	url := fmt.Sprintf("%s/library/book/isbn/%d", c.BaseURL, isbn)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return b, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.HTTPClient.Do(r)
	if err != nil {
		return b, fmt.Errorf("making request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return b, ErrBookNotFound
		}
		return b, ErrUnexpectedResponse
	}

	if err = json.NewDecoder(resp.Body).Decode(&b); err != nil {
		return b, ErrInvalidResponse
	}

	return b, nil
}
