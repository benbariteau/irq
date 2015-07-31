package view

import (
	"github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "strconv"
    "github.com/firba1/irq/model"
    "fmt"
)

func Quote(r render.Render, params martini.Params) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		env := map[string]interface{}{
			"title": "error",
			"error": "invalid quote id",
		}
		r.HTML(404, "error", env)
		return
	}

	db, err := model.NewModel("quotes.db")
	if err != nil {
		env := map[string]interface{}{
			"title": "error",
			"error": "db connection failed",
		}
		r.HTML(500, "error", env)
		return
	}

	quote, err := db.GetQuote(id)
	if err != nil {
		env := map[string]interface{}{
			"title": "error",
			"error": "quote not found",
		}
		r.HTML(404, "error", env)
		return
	}

	env := map[string]interface{}{
		"title": fmt.Sprintf("#%d", quote.ID),
		"quotes": []model.Quote{quote},
        "showPagination": false,
	}
	r.HTML(200, "quote", env)
}
