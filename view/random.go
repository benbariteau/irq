package view

import (
    "github.com/martini-contrib/render"
    "github.com/firba1/irq/model"
    "fmt"
)

func Random(r render.Render) {
	db, err := model.NewModel("quotes.db")
	if err != nil {
		env := map[string]interface{}{
			"title": "error",
			"error": "db connection failed",
		}
		r.HTML(500, "error", env)
		return
	}

	quotes, err := db.GetQuotes(model.Query{
        Limit:   1,
        OrderBy: []string{"random()"},
    })
    r.Redirect(fmt.Sprintf("/quote/%d", quotes[0].ID))
}
