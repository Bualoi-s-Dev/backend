package dto

type Pagination struct {
	Page    int `json:"page"`
	Limit   int `json:"limit"`
	MaxPage int `json:"maxPage"`
	Total   int `json:"total"`
}
