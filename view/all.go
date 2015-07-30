package view

import (
    "github.com/martini-contrib/render"
    "github.com/firba1/irq/model"
)

func All(r render.Render) {
	db, err := model.NewModel("quotes.db")
	if err != nil {
		env := map[string]interface{}{
			"title": "error",
			"error": "db connection failed",
		}
		r.HTML(500, "error", env)
		return
	}

	allQuotes, err := db.GetQuotes(0, 0)
	if err != nil {
		env := map[string]interface{}{
			"title": "error",
			"error": "failed to get quotes",
		}
		r.HTML(404, "error", env)
		return
	}

	env := map[string]interface{}{
		"title": "all quotes",
		"quotes": allQuotes,
	}
	r.HTML(200, "quote", env)
}
