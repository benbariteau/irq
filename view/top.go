package view

import (
	"github.com/firba1/irq/model"
	"github.com/martini-contrib/render"
	"net/http"
)

func Top(db model.Model, r render.Render, req *http.Request) {
	QuotesBase(db, r, req, "Top", []string{"score DESC"})
}
