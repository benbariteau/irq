package view

import (
	"net/http"

	"github.com/firba1/irq/model"

	"github.com/martini-contrib/render"
)

func Submit(r render.Render) {
	r.HTML(200, "submit", PageEnv{Title: "Submit"})
}

func SubmitForm(db model.Model, r render.Render, request *http.Request) {
	text := request.FormValue("text")
	isOffensive := request.FormValue("is_offensive")
	isNishbot := request.FormValue("is_nishbot")

	err := db.AddQuote(model.Quote{
		Text:        text,
		IsOffensive: isOffensive == "on",
		IsNishbot:   isNishbot == "on",
	})
	if err != nil {
		env := ErrorPageEnv{
			PageEnv{Title: "error"},
			ErrorEnv{ErrorMessage: "unable to add quote"},
		}
		r.HTML(500, "error", env)
	}

	r.Redirect("/latest")
}
