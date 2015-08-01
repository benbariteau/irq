package main

import (
	"net/http"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/firba1/irq/model"
	"github.com/firba1/irq/view"
)

func main() {
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Layout: "base",
	}))

	m.Use(func(req *http.Request, c martini.Context) {
		c.Map(view.IsJson(strings.HasSuffix(req.URL.Path, ".json")))
	})

	// middleware to inject DB connection into each request
	m.Use(func(r render.Render, c martini.Context, isJson view.IsJson) {
		db, err := model.NewModel("quotes.db")
		if err != nil {
			view.RenderError(r, 500, isJson, "db connection failed")
			return
		}
		c.Map(db)
	})

	m.Get("/", view.Index)
	m.Get("/latest(.json|)", view.Latest)
	m.Get("/all(.json|)", view.All)
	m.Get("/random(.json|)", view.Random)
	m.Get("/search(.json|)", view.Search)
	m.Get("/submit", view.Submit)
	m.Get("/top(.json|)", view.Top)
	m.Get("/quote/:id", view.Quote)

	m.NotFound(func(r render.Render, isJson view.IsJson) {
		view.RenderError(r, 404, isJson, "not found")
		return
	})

	m.Post("/submit", view.SubmitForm)
	m.Post("/vote/:id", view.Vote)
	m.Delete("/quote/:id", view.Delete)

	m.Run()
}
