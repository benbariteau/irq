package view

import (
	"github.com/firba1/irq/model"
)

type PageEnv struct {
	Title string
}

type ErrorEnv struct {
	ErrorMessage string `json:"error_message"`
}

type ErrorPageEnv struct {
	PageEnv
	ErrorEnv
}

type quoteEnv struct {
	PageEnv
	Quotes         []model.Quote
	ShowPagination bool
	Count          int
	Page           int
	PreviousPage   int
	NextPage       int
	Total          int
	MaxPage        int
	Query          string
}
