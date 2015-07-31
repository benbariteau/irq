package view

import (
	"net/http"

	"github.com/martini-contrib/render"

	"github.com/firba1/irq/model"
)

func Latest(db model.Model, r render.Render, req *http.Request) {
	QuotesBase(db, r, req, "Latest", []string{"id DESC"})
}
