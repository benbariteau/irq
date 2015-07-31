package view

import (
	"github.com/martini-contrib/render"
)

func About(r render.Render) {
	env := map[string]interface{}{
		"title": "About",
	}
	r.HTML(200, "about", env)
}
