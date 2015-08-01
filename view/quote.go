package view

import (
	"fmt"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/firba1/irq/model"
)

func Quote(db model.Model, r render.Render, params martini.Params) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		env := ErrorPageEnv{
			PageEnv{Title: "error"},
			ErrorEnv{ErrorMessage: "invalid quote id"},
		}
		r.HTML(404, "error", env)
		return
	}

	quote, err := db.GetQuote(id)
	if err != nil {
		env := ErrorPageEnv{
			PageEnv{Title: "error"},
			ErrorEnv{ErrorMessage: "quote not found"},
		}
		r.HTML(404, "error", env)
		return
	}

	env := quotePageEnv{
		PageEnv:        PageEnv{Title: fmt.Sprintf("#%d", quote.ID)},
		Quotes:         []model.Quote{quote},
		ShowPagination: false,
	}
	r.HTML(200, "quote", env)
}
