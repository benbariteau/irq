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

	err := db.AddQuote(model.Quote{
		Text:        text,
		IsOffensive: isOffensive == "on",
	})
	if err != nil {
		RenderError(r, 404, IsJson(false), "unable to add quote")
	}

	r.Redirect("/latest")
}
