package view

import (
	"github.com/martini-contrib/render"
	"net/http"
)

func Latest(r render.Render, req *http.Request) {
    QuotesBase(r, req, "Latest", []string{"id DESC"})
}
