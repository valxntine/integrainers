package library

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valxntine/integrainers/entity"
	"net/http"
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

func (c Client) GetBook(ctx context.Context, iban string) (entity.Book, error) {
	b := entity.Book{}

	url := fmt.Sprintf("%s/library/book/iban/%s", c.BaseURL, iban)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return b, errors.New("error creating request")
	}

	resp, err := c.HTTPClient.Do(r)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return b, errors.New("book not found")
		}
		return b, errors.New("failed to respond")
	}

	if err = json.NewDecoder(resp.Body).Decode(&b); err != nil {
		return b, errors.New("invalid json body")
	}

	return b, nil
}
