package view

import (
	"github.com/martini-contrib/render"
	"net/http"
)

func All(r render.Render, req *http.Request) {
    QuotesBase(r, req, "All", []string{"id ASC"})
}
