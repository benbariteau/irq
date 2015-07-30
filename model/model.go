package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Model struct {
	db *sql.DB
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

func (m Model) Close() error {
	return m.db.Close()
}
