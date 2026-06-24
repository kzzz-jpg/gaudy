package model

type Gua struct {
	GuaId     int      `json:"gua_id"`
	Title     string   `json:"title"`
	People    []string `json:"people"`
	Content   string   `json:"content"`
	CreatedAt string   `json:"created_at"`
}
