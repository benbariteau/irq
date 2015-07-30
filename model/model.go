package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Model struct {
	db *sql.DB
}

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

/*
NewModel creates a new model with a DB connection to the give dbPath (currently sqlite)
*/
func NewModel(dbPath string) (m Model, err error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return
	}
	m.db = db
	return
}

func (m Model) GetQuote(id int) (quote Quote, err error) {
	rawQ := rawQuote{}

	err = m.db.QueryRow(
		"SELECT id, text, score, time_created, is_offensive, is_nishbot from quote where id = ?",
		id,
	).Scan(
		&rawQ.ID,
		&rawQ.Text,
		&rawQ.Score,
		&rawQ.TimeCreated,
		&rawQ.IsOffensive,
		&rawQ.IsNishbot,
	)
	if err != nil {
		return
	}

	return toQuote(rawQ), nil
}

func toQuote(rawQ rawQuote) Quote {
	return Quote{
		ID:          rawQ.ID,
		Text:        rawQ.Text,
		Score:       rawQ.Score,
		TimeCreated: time.Unix(rawQ.TimeCreated, 0),
		IsOffensive: rawQ.IsOffensive != 0,
		IsNishbot:   rawQ.IsNishbot != 0,
	}
}

func (m Model) GetQuotes(limit, offset int, orderBy ...string) (quotes []Quote, err error) {
	rows, err := m.db.Query(
		strings.Join(
			[]string{
				"SELECT id, text, score, time_created, is_offensive, is_nishbot",
				"FROM quote",
				genOrderBy(orderBy),
				genLimitOffsetClause(limit, offset),
			},
			" ",
		),
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		rawQ := rawQuote{}
		err = rows.Scan(
			&rawQ.ID,
			&rawQ.Text,
			&rawQ.Score,
			&rawQ.TimeCreated,
			&rawQ.IsOffensive,
			&rawQ.IsNishbot,
		)
		if err != nil {
			return
		}
		quotes = append(quotes, toQuote(rawQ))
	}
	return
}

func genOrderBy(orderByColumns []string) string {
	if len(orderByColumns) == 0 {
		return ""
	}
	return "ORDER BY " + strings.Join(orderByColumns, ", ")
}

func genLimitOffsetClause(limit, offset int) string {
	if limit == 0 {
		return ""
	}
	return fmt.Sprint("LIMIT ", limit, " ", genOffsetClause(offset))
}

func genOffsetClause(offset int) string {
	if offset == 0 {
		return ""
	}
	return fmt.Sprint("OFFSET ", offset)
}

func (m Model) Close() error {
	return m.db.Close()
}
