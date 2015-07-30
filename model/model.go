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

func (m Model) AddQuote(q Quote) (err error) {
	rawQ := fromQuote(q)
	_, err = m.db.Exec(
		"INSERT INTO quote(text, score, time_created, is_offensive, is_nishbot) values(?, ?, ?, ?, ?)",
		rawQ.Text,
		rawQ.Score,
		time.Now().Unix(),
		rawQ.IsOffensive,
		rawQ.IsNishbot,
	)
	return
}

func fromQuote(quote Quote) rawQuote {
	return rawQuote{
		ID:          quote.ID,
		Text:        quote.Text,
		Score:       quote.Score,
		TimeCreated: quote.TimeCreated.Unix(),
		IsOffensive: boolToInt(quote.IsOffensive),
		IsNishbot:   boolToInt(quote.IsNishbot),
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func (m Model) DeleteQuote(id int) (err error) {
	_, err = m.db.Exec("DELETE FROM quote WHERE id = ?", id)
	return
}

func (m Model) VoteQuote(id int, vote int) (err error) {
	q, err := m.GetQuote(id)
	if err != nil {
		return
	}

	newScore := q.Score + vote

	_, err = m.db.Exec("UPDATE quote SET score = ? where id = ?", newScore, id)
	return
}

func voteRuneToInt(vote rune) int {
	switch vote {
	case '+':
		return 1
	case '-':
		return -1
	default:
		return 0
	}
}

func (m Model) Close() error {
	return m.db.Close()
}
