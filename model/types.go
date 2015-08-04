package model

import (
	"fmt"
	"strings"
	"time"
)

type Quote struct {
	ID          int       `json:"id"`
	Text        string    `json:"text"`
	Score       int       `json:"score"`
	TimeCreated time.Time `json:"time_created"`
	IsOffensive bool      `json:"is_offensive"`
	IsNishbot   bool      `json:"is_nishbot"`
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
	Search   string
	Limit    int
	Offset   int
	OrderBy  []string
	MaxLines int
	Clean    bool
}

func (q Query) WhereClause() string {
	parts := make([]string, 0, 3)
	if q.Search != "" {
		parts = append(parts, "text LIKE ?")
	}
	if q.MaxLines != 0 {
		parts = append(
			parts,
			fmt.Sprint("LENGTH(text) - LENGTH(REPLACE(text, X'0A', '')) + 1 <= ", q.MaxLines),
		)
	}
	if q.Clean {
		parts = append(parts, "not is_offensive and not is_nishbot")
	}
	if len(parts) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(parts, " AND ")
}

func (q Query) toSQL() string {
	parts := make([]string, 0, 3)

	whereClause := q.WhereClause()
	if whereClause != "" {
		parts = append(parts, whereClause)
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
