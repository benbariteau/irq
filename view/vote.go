package view

import (
	"github.com/firba1/irq/model"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
)

func Vote(db model.Model, req *http.Request, r render.Render, params martini.Params) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		r.JSON(404, ErrorEnv{"invalid quote id"})
		return
	}

	count, err := strconv.Atoi(req.FormValue("count"))
	if err != nil {
		r.JSON(404, ErrorEnv{"invalid vote count"})
		return
	}

	err = db.VoteQuote(id, count)
	if err != nil {
		r.JSON(500, ErrorEnv{"unable to vote quote"})
		return
	}

	quote, err := db.GetQuote(id)
	if err != nil {
		r.JSON(404, ErrorEnv{"quote not found"})
		return
	}

	env := struct {
		Score int `json:"score"`
	}{quote.Score}
	r.JSON(200, env)
}
