package view

import (
	"github.com/firba1/irq/model"
	"github.com/martini-contrib/render"
	"net/http"
)

func All(db model.Model, r render.Render, req *http.Request) {
	QuotesBase(db, r, req, "All", []string{"id ASC"})
}
