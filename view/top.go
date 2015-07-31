package view

import (
	"github.com/martini-contrib/render"
	"net/http"
)

func Top(r render.Render, req *http.Request) {
    QuotesBase(r, req, "Top", []string{"score DESC"})
}
