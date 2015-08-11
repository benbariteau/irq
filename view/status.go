package view

import (
	"github.com/martini-contrib/render"
)

func Status(r render.Render) {
	resp := struct {
		Status string `json:"status"`
	}{"up"}
	r.JSON(200, resp)
}
