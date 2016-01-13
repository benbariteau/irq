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
	Tags        []string  `json:"tags"`
}

type rawQuote struct {
	ID          int
	Text        string
	Score       int
	TimeCreated interface{}
	Tags        string
}

type Query struct {
	Search      string
	Limit       int
	Offset      int
	OrderBy     []string
	MaxLines    int
	IncludeTags []string
	ExcludeTags []string
}

func (q Query) WhereClause() string {
	parts := make([]string, 0, 3)
	if q.Search != "" {
		parts = append(parts, "quote.text LIKE ?")
	}
	if q.MaxLines != 0 {
		parts = append(
			parts,
			fmt.Sprint("LENGTH(quote.text) - LENGTH(REPLACE(quote.text, X'0A', '')) + 1 <= ", q.MaxLines),
		)
	}
	if len(q.IncludeTags) != 0 {
		parts = append(
			parts,
			fmt.Sprintf(
				`quote.id IN (SELECT quote_id FROM quote_tag WHERE tag IN (%s) GROUP BY quote_id HAVING count(quote_id) = %d)`,
				nArgs(len(q.IncludeTags)),
				len(q.IncludeTags),
			),
		)
	}
	if len(q.ExcludeTags) != 0 {
		parts = append(
			parts,
			fmt.Sprintf(
				`quote.id NOT IN (SELECT quote_id FROM quote_tag WHERE tag IN (%s))`,
				nArgs(len(q.ExcludeTags)),
			),
		)
	}
	if len(parts) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(parts, " AND ")
}

func nArgs(n int) string {
	return strings.Join(strings.Split(strings.Repeat("?", n), ""), ",")
}

func (q Query) toSQL() string {
	parts := make([]string, 0, 3)

	whereClause := q.WhereClause()
	if whereClause != "" {
		parts = append(parts, whereClause)
	}

	parts = append(parts, "GROUP BY quote.id")

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
