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
	db, err := sql.Open(dbType, configureDbPath(dbType, dbPath))
	if err != nil {
		return
	}
	m.db = db
	return
}

func configureDbPath(dbType, dbPath string) string {
	if dbType == "mysql" {
		dbPath = dbPath + "?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"
	}
	return dbPath
}

func (m Model) GetQuote(id int) (quote Quote, err error) {
	rawQ := rawQuote{}

	query := `
        SELECT quote.id, quote.text, quote.score, quote.time_created, ifnull(group_concat(tag), '')
        FROM quote
        LEFT JOIN quote_tag ON (quote_tag.quote_id = quote.id)
        WHERE quote.id = ?
        GROUP BY quote.id
    `

	err = m.db.QueryRow(
		query,
		id,
	).Scan(
		&rawQ.ID,
		&rawQ.Text,
		&rawQ.Score,
		&rawQ.TimeCreated,
		&rawQ.Tags,
	)
	if err != nil {
		return
	}

	return toQuote(rawQ), nil
}

func (m Model) GetPopularTags(n int) (tags []string, err error) {
	query := `SELECT TAG FROM quote_tag GROUP BY tag ORDER BY COUNT(*) DESC LIMIT ?`

	rows, err := m.db.Query(query, n)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return
		}
		tags = append(tags, tag)
	}
	return
}

func (m Model) GetQuotes(q Query) (quotes []Quote, err error) {
	query := `
        SELECT quote.id, quote.text, quote.score, quote.time_created, ifnull(group_concat(tag), '')
        FROM quote
        LEFT JOIN quote_tag ON (quote_tag.quote_id = quote.id)
    ` + "\n" + q.toSQL()

	queryArgs := make([]interface{}, 0, 1)
	if q.Search != "" {
		queryArgs = append(queryArgs, "%"+q.Search+"%")
	}
	for _, tag := range q.IncludeTags {
		queryArgs = append(queryArgs, tag)
	}
	for _, tag := range q.ExcludeTags {
		queryArgs = append(queryArgs, tag)
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
			&rawQ.Tags,
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
	for _, tag := range q.IncludeTags {
		queryArgs = append(queryArgs, tag)
	}
	for _, tag := range q.ExcludeTags {
		queryArgs = append(queryArgs, tag)
	}

	err = m.db.QueryRow(query, queryArgs...).Scan(&count)
	if err != nil && err == sql.ErrNoRows {
		err = nil
		count = 0
	}
	return
}

func (m Model) CountAllQuotes() (count int, err error) {
	return m.CountQuotes(Query{})
}

func (m Model) AddQuote(q Quote) (err error) {
	rawQ := fromQuote(q)
	result, err := m.db.Exec(
		"INSERT INTO quote(text, score, time_created, is_offensive, is_nishbot) values(?, ?, ?, ?, ?)",
		rawQ.Text,
		rawQ.Score,
		time.Now().Unix(),
		false,
		false,
	)
	if err != nil {
		return
	}
	quote_id, err := result.LastInsertId()
	if err != nil {
		return
	}
	for _, tag := range q.Tags {
		m.db.Exec(
			"INSERT INTO quote_tag(quote_id, tag) values (?, ?)",
			quote_id,
			tag,
		)
	}
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
		false,
		false,
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
