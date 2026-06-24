package model

type Gua struct {
	GuaId   int
	Title   string   `json:"title"`
	People  []string `json:"people"`
	Content string   `json:"content"`
}
