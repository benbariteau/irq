package view

import (
	"github.com/firba1/irq/model"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
)

func Vote(req *http.Request, r render.Render, params martini.Params) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		r.JSON(404, errorEnv{"invalid quote id"})
		return
	}

	count, err := strconv.Atoi(req.FormValue("count"))
	if err != nil {
		r.JSON(404, errorEnv{"invalid vote count"})
		return
	}

	db, err := model.NewModel("quotes.db")
	if err != nil {
		r.JSON(500, errorEnv{"what happen db conn no work"})
		return
	}

	err = db.VoteQuote(id, count)
	if err != nil {
		r.JSON(500, errorEnv{"unable to vote quote"})
		return
	}

	quote, err := db.GetQuote(id)
	if err != nil {
		r.JSON(404, errorEnv{"quote not found"})
		return
	}

	env := struct {
		Score int `json:"score"`
	}{quote.Score}
	r.JSON(200, env)
}
