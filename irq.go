package main

import (
	"github.com/firba1/irq/view"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func main() {
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Layout: "base",
	}))

	m.Get("/", view.Index)
	m.Get("/latest", view.Latest)
	m.Get("/all", view.All)
	m.Get("/about", view.About)
	m.Get("/random", view.Random)
	m.Get("/random.json", view.RandomJson)
	m.Get("/top", view.Top)
	m.Get("/search", view.Search)
	m.Get("/submit", view.Submit)
	m.Get("/quote/:id", view.Quote)

	m.Post("/submit", view.SubmitForm)
	m.Post("/vote/:id", view.Vote)

	m.Run()
}
