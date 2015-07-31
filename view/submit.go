package view

import (
	"net/http"

	"github.com/firba1/irq/model"

	"github.com/martini-contrib/render"
)

func Submit(r render.Render) {
	env := map[string]interface{}{
		"title": "Submit",
	}
	r.HTML(200, "submit", env)
}

func SubmitForm(r render.Render, request *http.Request) {
	text := request.FormValue("text")
	isOffensive := request.FormValue("is_offensive")
	isNishbot := request.FormValue("is_nishbot")

	db, err := model.NewModel("quotes.db")
	if err != nil {
		env := map[string]interface{}{
			"title": "error",
			"error": "what happen db conn no work",
		}
		r.HTML(500, "error", env)
	}

	err = db.AddQuote(model.Quote{
		Text:        text,
		IsOffensive: isOffensive == "on",
		IsNishbot:   isNishbot == "on",
	})
	if err != nil {
		env := map[string]interface{}{
			"title": "error",
			"error": "unable to add quote",
		}
		r.HTML(500, "error", env)
	}

	r.Redirect("/latest")
}
