package view

import (
	"net/http"

	"github.com/firba1/irq/model"

	"github.com/martini-contrib/render"
)

func Submit(r render.Render) {
	r.HTML(200, "submit", pageEnv{Title: "Submit"})
}

func SubmitForm(r render.Render, request *http.Request) {
	text := request.FormValue("text")
	isOffensive := request.FormValue("is_offensive")
	isNishbot := request.FormValue("is_nishbot")

	db, err := model.NewModel("quotes.db")
	if err != nil {
		env := errorPageEnv{
			pageEnv{Title: "error"},
			errorEnv{ErrorMessage: "what happen db conn no work"},
		}
		r.HTML(500, "error", env)
		return
	}

	err = db.AddQuote(model.Quote{
		Text:        text,
		IsOffensive: isOffensive == "on",
		IsNishbot:   isNishbot == "on",
	})
	if err != nil {
		env := errorPageEnv{
			pageEnv{Title: "error"},
			errorEnv{ErrorMessage: "unable to add quote"},
		}
		r.HTML(500, "error", env)
	}

	r.Redirect("/latest")
}
