package model

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type Model struct {
	db *sql.DB
}

/*
NewModel creates a new model with a DB connection to the give dbPath (currently sqlite)
*/
func NewModel(dbType, dbPath string) (m Model, err error) {
	db, err := sql.Open(dbType, dbPath)
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

func (m Model) GetQuotes(q Query) (quotes []Quote, err error) {
	query := strings.Join(
		[]string{
			"SELECT id, text, score, time_created, is_offensive, is_nishbot",
			"FROM quote",
			q.toSQL(),
		},
		"\n",
	)
	queryArgs := make([]interface{}, 0, 1)
	if q.Search != "" {
		queryArgs = append(queryArgs, "%"+q.Search+"%")
	}

	rows, err := m.db.Query(query, queryArgs...)
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

func (m Model) CountQuotes(q Query) (count int, err error) {
	queryParts := []string{
		"SELECT count(*)",
		"FROM quote",
	}
	whereClause := q.WhereClause()
	if whereClause != "" {
		queryParts = append(queryParts, whereClause)
	}
	query := strings.Join(queryParts, "\n")

	queryArgs := make([]interface{}, 0, 1)
	if q.Search != "" {
		queryArgs = append(queryArgs, "%"+q.Search+"%")
	}

	err = m.db.QueryRow(query, queryArgs...).Scan(&count)
	return
}

func (m Model) CountAllQuotes() (count int, err error) {
	return m.CountQuotes(Query{})
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

func (m Model) DeleteQuote(id int) (err error) {
	quote, err := m.GetQuote(id)
	if err != nil {
		return
	}
	_, err = m.db.Exec("DELETE FROM quote WHERE id = ?", id)
	if err != nil {
		return
	}

	rawQ := fromQuote(quote)
	_, err = m.db.Exec(
		"INSERT INTO deleted_quote(text, score, time_created, is_offensive, is_nishbot) values(?, ?, ?, ?, ?)",
		rawQ.Text,
		rawQ.Score,
		rawQ.TimeCreated,
		rawQ.IsOffensive,
		rawQ.IsNishbot,
	)
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

func (m Model) Close() error {
	return m.db.Close()
}
