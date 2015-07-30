package model

import (
	"database/sql"
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
	TimeCreated int
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

	return Quote{
		ID:          rawQ.ID,
		Text:        rawQ.Text,
		Score:       rawQ.Score,
		TimeCreated: time.Unix(int64(rawQ.TimeCreated), 0),
		IsOffensive: rawQ.IsOffensive != 0,
		IsNishbot:   rawQ.IsNishbot != 0,
	}, nil
}

func (m Model) Close() error {
	return m.db.Close()
}
