package model

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestNewModel(t *testing.T) {
	model, err := NewModel("./test.db")
	if err != nil {
		t.Error("Got non-nil error: ", err)
	}
	if model.db == nil {
		t.Error("Got nil db connection")
	}
}

type TestModel struct {
	m      Model
	dbfile *os.File
}

func (tm TestModel) Close() {
	tm.m.Close()
	if tm.dbfile != nil {
		dbfilepath, err := filepath.Abs(filepath.Dir(tm.dbfile.Name()))
		if err != nil {
			os.Remove(dbfilepath)
		}
	}
}

func createTestModel() (tm TestModel, err error) {
	// create temp sqlite DB
	f, err := ioutil.TempFile("", "quotedb")
	if err != nil {
		return
	}
	tm.dbfile = f

	// read schema file
	schemaBytes, err := ioutil.ReadFile("../schema/quote.sql")
	if err != nil {
		return
	}
	schema := string(schemaBytes)

	// open db conn
	m, err := NewModel(f.Name())
	if err != nil {
		return
	}
	tm.m = m

	// create table(s)
	_, err = m.db.Exec(schema)
	if err != nil {
		return
	}

	_, err = m.db.Exec(
		"insert into quote(id, text, score, time_created, is_offensive, is_nishbot) values(?, ?, ?, ?, ?, ?)",
		1,
		"fart joke",
		0,
		0,
		0,
		0,
	)
	if err != nil {
		return
	}

	return
}

func TestGetQuote(t *testing.T) {
	tm, err := createTestModel()
	defer tm.Close()
	if err != nil {
		t.Error("Got unexpected error: ", err)
	}

	quote, err := tm.m.GetQuote(1)
	if err != nil {
		t.Error("Got non-nil error: ", err)
	}

	expected := Quote{
		ID:          1,
		Text:        "fart joke",
		Score:       0,
		TimeCreated: time.Unix(0, 0),
		IsOffensive: false,
		IsNishbot:   false,
	}

	if !reflect.DeepEqual(quote, expected) {
		t.Error("Got: ", quote, "\nExpected:", expected)
	}
}
