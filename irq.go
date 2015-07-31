package main

import (
	"reflect"

	"github.com/firba1/irq/model"
	"github.com/firba1/irq/view"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func main() {
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Layout: "base",
	}))

	// middleware to inject DB connection into each request
	m.Use(func(r render.Render, c martini.Context) {
		db, err := model.NewModel("quotes.db")
		if err != nil {
			env := view.ErrorPageEnv{
				view.PageEnv{Title: "error"},
				view.ErrorEnv{ErrorMessage: "db connection failed"},
			}
			r.HTML(500, "error", env)
			return
		}
		c.Set(reflect.TypeOf(db), reflect.ValueOf(db))
	})

	m.Get("/", view.Index)
	m.Get("/latest", view.Latest)
	m.Get("/all", view.All)
	m.Get("/random", view.Random)
	m.Get("/random.json", view.RandomJson)
	m.Get("/top", view.Top)
	m.Get("/search", view.Search)
	m.Get("/submit", view.Submit)
	m.Get("/quote/:id", view.Quote)

	m.NotFound(func(r render.Render) {
		env := view.ErrorPageEnv{
			view.PageEnv{Title: "error"},
			view.ErrorEnv{ErrorMessage: "page not found"},
		}
		r.HTML(404, "error", env)
		return
	})

	m.Post("/submit", view.SubmitForm)
	m.Post("/vote/:id", view.Vote)
	m.Delete("/quote/:id", view.Delete)

	m.Run()
}
