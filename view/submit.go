package view

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/firba1/irq/model"

	"github.com/martini-contrib/render"
)

func Submit(db model.Model, r render.Render) {
	popularTags, err := db.GetPopularTags(10)
	if err != nil {
		RenderError(r, 404, IsJson(false), fmt.Sprint("unable to get tags", err))
		return
	}

	r.HTML(200, "submit", submitPageEnv{
		PageEnv: PageEnv{
			Title: "Submit",
		},
		PopularTags: popularTags,
	})
}

func SubmitForm(db model.Model, r render.Render, request *http.Request) {
	text := request.FormValue("text")
	tags := strings.Split(strings.ToLower(request.FormValue("tags")), ",")

	err := db.AddQuote(model.Quote{
		Text: text,
		Tags: tags,
	})
	if err != nil {
		RenderError(r, 404, IsJson(false), fmt.Sprint("unable to add quote", err))
	}

	r.Redirect("/latest")
}
