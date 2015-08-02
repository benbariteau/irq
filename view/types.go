package view

import (
	"github.com/firba1/irq/model"
)

type IsJson bool

type PageEnv struct {
	Title string
	Query string
}

type ErrorEnv struct {
	ErrorMessage string `json:"error_message"`
}

type ErrorPageEnv struct {
	PageEnv
	ErrorEnv
}

type quotePageEnv struct {
	PageEnv
	Quotes          []model.Quote
	ShowPagination  bool
	Count           int
	Page            int
	PreviousPageURL string
	NextPageURL     string
	Total           int
	MaxPage         int
}
