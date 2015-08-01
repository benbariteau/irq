package view

import (
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/firba1/irq/model"
)

func maxPage(totalItems, perPage int) int {
	return (totalItems-1)/perPage + 1
}

func QuotesBase(title string, orderBy []string) martini.Handler {
	return func(db model.Model, r render.Render, req *http.Request, isJson IsJson) {
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
			RenderError(r, 404, isJson, "failed to get quotes")
			return
		}

		total, err := db.CountQuotes(query)
		if err != nil {
			RenderError(r, 404, isJson, "failed to get quotes")
			return
		}

		if isJson {
			env := struct {
				Quotes []model.Quote `json:"quotes"`
				Total  int           `json:"total"`
			}{quotes, total}
			r.JSON(200, env)
			return
		}

		maxPage := maxPage(total, count)
		previousPage := page - 1
		nextPage := page + 1
		if nextPage > maxPage {
			nextPage = 0
		}

		env := quotePageEnv{
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
}

func RenderError(r render.Render, code int, isJson IsJson, errorMessage string) {
	env := ErrorEnv{ErrorMessage: errorMessage}
	if isJson {
		r.JSON(code, env)
	} else {
		env := ErrorPageEnv{
			PageEnv{Title: "error"},
			env,
		}
		r.HTML(code, "error", env)
	}
}
