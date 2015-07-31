package view

import (
	"github.com/firba1/irq/model"
	"github.com/martini-contrib/render"
	"net/http"
)

func Search(db model.Model, r render.Render, req *http.Request) {
	QuotesBase(db, r, req, "Search", []string{"score DESC"})
}
