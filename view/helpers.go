package view

import (
	"net/http"
	"strconv"

	"github.com/martini-contrib/render"

	"github.com/firba1/irq/model"
)

func maxPage(totalItems, perPage int) int {
	return (totalItems-1)/perPage + 1
}

func QuotesBase(db model.Model, r render.Render, req *http.Request, title string, orderBy []string) {
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

	offset := (page - 1) * count
	quotes, err := db.GetQuotes(model.Query{
		Limit:   count,
		Offset:  offset,
		Search:  query,
		OrderBy: orderBy,
	})
	if err != nil {
		env := ErrorPageEnv{
			PageEnv{Title: "error"},
			ErrorEnv{ErrorMessage: "failed to get quotes"},
		}
		r.HTML(404, "error", env)
		return
	}

	total, err := db.CountQuotes(query)
	if err != nil {
		env := ErrorPageEnv{
			PageEnv{Title: "error"},
			ErrorEnv{ErrorMessage: "failed to get quotes"},
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

	env := quoteEnv{
		PageEnv: PageEnv{
			Title: title,
		},
		Quotes:         quotes,
		ShowPagination: true,
		Count:          count,
		Page:           page,
		PreviousPage:   previousPage,
		NextPage:       nextPage,
		Total:          total,
		MaxPage:        maxPage,
		Query:          query,
	}
	r.HTML(200, "quote", env)
}
