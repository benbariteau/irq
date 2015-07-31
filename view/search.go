package view

import (
	"github.com/martini-contrib/render"
	"net/http"
)

func Search(r render.Render, req *http.Request) {
    QuotesBase(r, req, "Search", []string{"score DESC"})
}
