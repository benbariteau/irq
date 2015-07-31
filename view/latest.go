package view

import (
	"github.com/firba1/irq/model"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
)

func Latest(r render.Render, req *http.Request) {

	qs := req.URL.Query()

	page, err := strconv.Atoi(qs.Get("page"))
	if err != nil {
		page = 1
	}

	count, err := strconv.Atoi(qs.Get("count"))
	if err != nil || count == 0 {
		count = 20
	}

	db, err := model.NewModel("quotes.db")
	if err != nil {
		env := map[string]interface{}{
			"title": "error",
			"error": "db connection failed",
		}
		r.HTML(500, "error", env)
		return
	}

	offset := (page - 1) * count
	quotes, err := db.GetQuotes(model.Query{
		Limit:   count,
		Offset:  offset,
		OrderBy: []string{"id DESC"},
	})
	if err != nil {
		env := map[string]interface{}{
			"title": "error",
			"error": "failed to get quotes",
		}
		r.HTML(404, "error", env)
		return
	}

	allQuotes, err := db.GetQuotes(model.Query{})
	if err != nil {
		env := map[string]interface{}{
			"title": "error",
			"error": "failed to get quotes",
		}
		r.HTML(404, "error", env)
		return
	}

	total := len(allQuotes)
	maxPage := total/count + 1
	previousPage := page - 1
	nextPage := page + 1
	if nextPage > maxPage {
		nextPage = 0
	}

	env := map[string]interface{}{
		"title":          "Latest",
		"quotes":         quotes,
		"showPagination": true,
		"count":          count,
		"page":           page,
		"previousPage":   previousPage,
		"nextPage":       nextPage,
		"total":          total,
		"maxPage":        maxPage,
	}
	r.HTML(200, "quote", env)
}
