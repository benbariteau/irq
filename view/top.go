package view

import (
	"net/http"

	"github.com/martini-contrib/render"

	"github.com/firba1/irq/model"
)

func Top(db model.Model, r render.Render, req *http.Request) {
	QuotesBase(db, r, req, "Top", []string{"score DESC"})
}
