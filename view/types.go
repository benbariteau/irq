package view

import (
	"github.com/firba1/irq/model"
)

type pageEnv struct {
	Title string
}

type errorEnv struct {
	ErrorMessage string `json:"error_message"`
}

type errorPageEnv struct {
	pageEnv
	errorEnv
}

type quoteEnv struct {
	pageEnv
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
