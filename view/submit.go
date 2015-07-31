package view

import (
    "github.com/martini-contrib/render"
)

func Submit(r render.Render) {
	env := map[string]interface{}{
		"title": "Submit",
	}
	r.HTML(200, "submit", env)
}
