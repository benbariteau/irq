package model

import (
	"fmt"
	"strings"
	"time"
)

type Quote struct {
	ID          int
	Text        string
	Score       int
	TimeCreated time.Time
	IsOffensive bool
	IsNishbot   bool
}

type rawQuote struct {
	ID          int
	Text        string
	Score       int
	TimeCreated int64
	IsOffensive int
	IsNishbot   int
}

type Query struct {
	Search  string
	Limit   int
	Offset  int
	OrderBy []string
}

func (q Query) toSQL() string {
	parts := make([]string, 0, 4)

	if q.Search != "" {
		parts = append(parts, searchWhereClause(q.Search))
	}

	if len(q.OrderBy) != 0 {
		parts = append(parts, "ORDER BY "+strings.Join(q.OrderBy, ", "))
	}

	if q.Limit != 0 {
		parts = append(parts, fmt.Sprint("LIMIT ", q.Limit))
		if q.Offset != 0 {
			parts = append(parts, fmt.Sprint("OFFSET ", q.Offset))
		}
	}

	return strings.Join(parts, "\n")
}
