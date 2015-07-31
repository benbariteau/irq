package view

import (
	"fmt"
	"github.com/firba1/irq/model"
	"github.com/martini-contrib/render"
	"net/http"
)

func Random(db model.Model, r render.Render) {
	quotes, err := db.GetQuotes(model.Query{
		Limit:   1,
		OrderBy: []string{"random()"},
	})
	if err != nil || len(quotes) == 0 {
		env := ErrorPageEnv{
			PageEnv{Title: "error"},
			ErrorEnv{ErrorMessage: "quote not found"},
		}
		r.HTML(500, "error", env)
		return
	}

	r.Redirect(fmt.Sprintf("/quote/%d", quotes[0].ID))
}

func RandomJson(db model.Model, r render.Render, req *http.Request) {
	qs := req.URL.Query()

	query := qs.Get("query")

	quotes, err := db.GetQuotes(model.Query{
		Limit:   1,
		Search:  query,
		OrderBy: []string{"random()"},
	})

	if err != nil || len(quotes) == 0 {
		r.JSON(500, ErrorEnv{"quote not found"})
		return
	}

	quote := quotes[0]

	r.JSON(200, quote)
}
