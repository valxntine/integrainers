package entity

type Book struct {
	Author string `json:"author"`
	Name   string `json:"name"`
	ISBN   int    `json:"isbn"`
	Pages  int    `json:"pages"`
}
