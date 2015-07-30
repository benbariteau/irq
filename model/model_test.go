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

func createTestModel(quotes ...Quote) (tm TestModel, err error) {
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

	for _, quote := range quotes {
		rawQ := fromQuote(quote)
		_, err = m.db.Exec(
			"insert into quote(id, text, score, time_created, is_offensive, is_nishbot) values(?, ?, ?, ?, ?, ?)",
			rawQ.ID,
			rawQ.Text,
			rawQ.Score,
			rawQ.TimeCreated,
			rawQ.IsOffensive,
			rawQ.IsNishbot,
		)
		if err != nil {
			return
		}
	}

	return
}

func TestGetQuote(t *testing.T) {
	expected := Quote{
		ID:          1,
		Text:        "fart joke",
		Score:       0,
		TimeCreated: time.Unix(0, 0),
		IsOffensive: false,
		IsNishbot:   false,
	}
	tm, err := createTestModel(expected)
	defer tm.Close()
	if err != nil {
		t.Error("Got unexpected error: ", err)
	}

	quote, err := tm.m.GetQuote(1)
	if err != nil {
		t.Error("Got non-nil error: ", err)
	}

	if !reflect.DeepEqual(quote, expected) {
		t.Error("Got: ", quote, "\nExpected:", expected)
	}
}

func TestGetQuotes(t *testing.T) {
	dbquotes := []Quote{
		Quote{
			ID:          1,
			Text:        "fart joke",
			Score:       0,
			TimeCreated: time.Unix(0, 0),
			IsOffensive: true,
			IsNishbot:   false,
		},
		Quote{
			ID:          2,
			Text:        "javascript joke",
			Score:       -5,
			TimeCreated: time.Unix(10, 0),
			IsOffensive: false,
			IsNishbot:   false,
		},
		Quote{
			ID:          3,
			Text:        "python joke",
			Score:       10,
			TimeCreated: time.Unix(10, 0),
			IsOffensive: false,
			IsNishbot:   false,
		},
	}

	tm, err := createTestModel(dbquotes...)
	defer tm.Close()
	if err != nil {
		t.Error("Got unexpected error: ", err)
	}

	tests := []struct {
		limit    int
		offset   int
		orderby  []string
		expected []Quote
	}{
		{0, 0, []string{}, dbquotes},
		// limit
		{2, 0, []string{}, dbquotes[:2]},
		// limit and offset
		{2, 2, []string{}, dbquotes[2:]},
		// order (no ordering, default ascending)
		{
			0,
			0,
			[]string{"Score"},
			[]Quote{dbquotes[1], dbquotes[0], dbquotes[2]},
		},
		// order descending
		{
			0,
			0,
			[]string{"Score DESC"},
			[]Quote{dbquotes[2], dbquotes[0], dbquotes[1]},
		},
		// order descending
		{
			0,
			0,
			[]string{"time_created ASC", "score DESC"},
			[]Quote{dbquotes[0], dbquotes[2], dbquotes[1]},
		},
		// limit and ordering
		{
			2,
			0,
			[]string{"Score DESC"},
			[]Quote{dbquotes[2], dbquotes[0]},
		},
		// limit, ordering, and offset
		{
			2,
			2,
			[]string{"Score DESC"},
			[]Quote{dbquotes[1]},
		},
	}

	for _, test := range tests {
		quotes, err := tm.m.GetQuotes(test.limit, test.offset, test.orderby...)
		if err != nil {
			t.Error("Got unexpected error: ", err)
		}

		if !reflect.DeepEqual(quotes, test.expected) {
			t.Error("Got: ", quotes, "\nExpected:", test.expected)
		}
	}
}

func TestAddQuote(t *testing.T) {
	tm, err := createTestModel()
	defer tm.Close()
	if err != nil {
		t.Error("Got unexpected error: ", err)
	}

	quote := Quote{
		ID:   1,
		Text: "php joke",
	}
	tm.m.AddQuote(quote)

	gotQuote, err := tm.m.GetQuote(1)
	if err != nil {
		t.Error("Got unexpected error: ", err)
	}

	if gotQuote.ID != quote.ID {
		t.Error("quote.ID\nGot: ", gotQuote.ID, "\nExpected:", quote.ID)
	}
	if gotQuote.Text != quote.Text {
		t.Error("Got: ", gotQuote.Text, "\nExpected:", quote.Text)
	}
	if gotQuote.Score != quote.Score {
		t.Error("Got: ", gotQuote.Score, "\nExpected:", quote.Score)
	}
	if gotQuote.IsOffensive != quote.IsOffensive {
		t.Error("Got: ", gotQuote.IsOffensive, "\nExpected:", quote.IsOffensive)
	}
	if gotQuote.IsNishbot != quote.IsNishbot {
		t.Error("Got: ", gotQuote.IsNishbot, "\nExpected:", quote.IsNishbot)
	}
}
