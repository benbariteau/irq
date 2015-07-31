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
		env := map[string]interface{}{
			"error": "invalid quote id",
		}
		r.JSON(404, env)
		return
	}

	count, err := strconv.Atoi(req.FormValue("count"))
	if err != nil {
		env := map[string]interface{}{
			"error": "invalid vote count",
		}
		r.JSON(404, env)
		return
	}

	db, err := model.NewModel("quotes.db")
	if err != nil {
		env := map[string]interface{}{
			"error": "what happen db conn no work",
		}
		r.JSON(500, env)
	}

	err = db.VoteQuote(id, count)
	if err != nil {
		env := map[string]interface{}{
			"error": "unable to vote quote",
		}
		r.JSON(500, env)
	}

	quote, err := db.GetQuote(id)
	if err != nil {
		env := map[string]interface{}{
			"error": "quote not found",
		}
		r.JSON(404, env)
		return
	}

	env := map[string]interface{}{
		"score": quote.Score,
	}
	r.JSON(200, env)
}
