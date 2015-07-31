package view

import (
	"github.com/firba1/irq/model"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
)

func maxPage(totalItems, perPage int) int {
	return (totalItems - 1)/perPage + 1
}

func QuotesBase(r render.Render, req *http.Request, title string, orderBy []string) {
	qs := req.URL.Query()

	query := qs.Get("query")

	page, err := strconv.Atoi(qs.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	count, err := strconv.Atoi(qs.Get("count"))
	if err != nil || count < 1 {
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
        Limit: count,
        Offset: offset,
        Search: query,
        OrderBy: orderBy,
    })
	if err != nil {
		env := map[string]interface{}{
			"title": "error",
			"error": "failed to get quotes",
		}
		r.HTML(404, "error", env)
		return
	}

	total, err := db.CountQuotes(query)
	if err != nil {
		env := map[string]interface{}{
			"title": "error",
			"error": "failed to get quotes",
		}
		r.HTML(404, "error", env)
		return
	}

	maxPage := maxPage(total, count)
	previousPage := page - 1
	nextPage := page + 1
	if nextPage > maxPage {
		nextPage = 0
	}

	env := map[string]interface{}{
		"title":          title,
		"quotes":         quotes,
		"showPagination": true,
		"count":          count,
		"page":           page,
		"previousPage":   previousPage,
		"nextPage":       nextPage,
		"total":          total,
		"maxPage":        maxPage,
		"query":          query,
	}
	r.HTML(200, "quote", env)
}

