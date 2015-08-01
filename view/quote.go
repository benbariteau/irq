package view

import (
	"fmt"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/firba1/irq/model"
)

func Quote(db model.Model, r render.Render, params martini.Params, isJson IsJson) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		RenderError(r, 404, isJson, "invalid quote id")
		return
	}

	quote, err := db.GetQuote(id)
	if err != nil {
		RenderError(r, 404, isJson, "quote not found")
		return
	}

	if isJson {
		r.JSON(200, quote)
		return
	}

	env := quotePageEnv{
		PageEnv:        PageEnv{Title: fmt.Sprintf("#%d", quote.ID)},
		Quotes:         []model.Quote{quote},
		ShowPagination: false,
	}
	r.HTML(200, "quote", env)
}
