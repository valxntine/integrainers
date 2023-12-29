package entity

type Book struct {
	Author string `json:"author"`
	Name   string `json:"name"`
	Iban   string `json:"iban"`
	Pages  int    `json:"pages"`
}
