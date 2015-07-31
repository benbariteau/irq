package view

import (
	"fmt"
	"github.com/firba1/irq/model"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"strconv"
)

func Quote(r render.Render, params martini.Params) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		env := errorPageEnv{
			pageEnv{Title: "error"},
			errorEnv{ErrorMessage: "invalid quote id"},
		}
		r.HTML(404, "error", env)
		return
	}

	db, err := model.NewModel("quotes.db")
	if err != nil {
		env := errorPageEnv{
			pageEnv{Title: "error"},
			errorEnv{ErrorMessage: "db connection failed"},
		}
		r.HTML(500, "error", env)
		return
	}

	quote, err := db.GetQuote(id)
	if err != nil {
		env := errorPageEnv{
			pageEnv{Title: "error"},
			errorEnv{ErrorMessage: "quote not found"},
		}
		r.HTML(404, "error", env)
		return
	}

	env := quoteEnv{
		pageEnv:        pageEnv{Title: fmt.Sprintf("#%d", quote.ID)},
		Quotes:         []model.Quote{quote},
		ShowPagination: false,
	}
	r.HTML(200, "quote", env)
}
