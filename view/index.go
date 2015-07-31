package view

import (
	"github.com/martini-contrib/render"
)

func Index(r render.Render) {
    r.Redirect("/top")
}
