package main

import (
	"github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "strconv"
    "github.com/firba1/irq/model"
    "fmt"
)

func main() {
	m := martini.Classic()

    m.Use(render.Renderer(render.Options{
        Layout: "base",
    }))


    m.Get("/quote/:id", func(r render.Render, params martini.Params) {
        id, err := strconv.Atoi(params["id"])
        if err != nil {
            env := map[string]interface{}{
                "title": "error",
                "error": "invalid quote id",
            }
            r.HTML(404, "error", env)
            return
        }

        m, err := model.NewModel("quotes.db")
        if err != nil {
            env := map[string]interface{}{
                "title": "error",
                "error": "db connection failed",
            }
            r.HTML(500, "error", env)
            return
        }

        quote, err := m.GetQuote(id)
        if err != nil {
            env := map[string]interface{}{
                "title": "error",
                "error": "quote not found",
            }
            r.HTML(404, "error", env)
            return
        }

        env := map[string]interface{}{
            "title": fmt.Sprintf("#%d", quote.ID),
            "id": quote.ID,
            "score": quote.Score,
            "text": quote.Text,
        }
        r.HTML(200, "quote", env)
	})

	m.Run()
}
