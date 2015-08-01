package view

import (
	"fmt"
	"net/http"

	"github.com/martini-contrib/render"

	"github.com/firba1/irq/model"
)

func Random(db model.Model, r render.Render, req *http.Request, isJson IsJson) {
	qs := req.URL.Query()

	query := qs.Get("query")

	quotes, err := db.GetQuotes(model.Query{
		Limit:   1,
		Search:  query,
		OrderBy: []string{"random()"},
	})

	if err != nil || len(quotes) == 0 {
		RenderError(r, 500, isJson, "quote not found")
		return
	}

	quote := quotes[0]

	if isJson {
		r.JSON(200, quote)
	} else {
		r.Redirect(fmt.Sprintf("/quote/%d", quote.ID))
	}
}
