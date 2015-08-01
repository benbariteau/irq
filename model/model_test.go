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

	// open db conn
	m, err := NewModel(f.Name())
	if err != nil {
		return
	}
	tm.m = m

	// read schema files
	filenames, err := filepath.Glob("./schema/*.sql")
	if err != nil {
		return
	}

	for _, filename := range filenames {
		var schemaBytes []byte
		schemaBytes, err = ioutil.ReadFile(filename)
		if err != nil {
			return
		}
		schema := string(schemaBytes)

		// create table(s)
		_, err = m.db.Exec(schema)
		if err != nil {
			return
		}
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
		query    Query
		expected []Quote
	}{
		{Query{}, dbquotes},
		// limit
		{Query{Limit: 2}, dbquotes[:2]},
		// limit and offset
		{Query{Limit: 2, Offset: 2}, dbquotes[2:]},
		// order (no ordering, default ascending)
		{
			Query{OrderBy: []string{"score"}},
			[]Quote{dbquotes[1], dbquotes[0], dbquotes[2]},
		},
		// order descending
		{
			Query{OrderBy: []string{"score DESC"}},
			[]Quote{dbquotes[2], dbquotes[0], dbquotes[1]},
		},
		// order descending
		{
			Query{OrderBy: []string{"time_created ASC", "score DESC"}},
			[]Quote{dbquotes[0], dbquotes[2], dbquotes[1]},
		},
		// limit and ordering
		{
			Query{
				Limit:   2,
				OrderBy: []string{"Score DESC"},
			},
			[]Quote{dbquotes[2], dbquotes[0]},
		},
		// limit, ordering, and offset
		{
			Query{
				Limit:   2,
				Offset:  2,
				OrderBy: []string{"Score DESC"},
			},
			[]Quote{dbquotes[1]},
		},
		// searching
		{
			Query{Search: "javascript"},
			[]Quote{dbquotes[1]},
		},
	}

	for _, test := range tests {
		quotes, err := tm.m.GetQuotes(test.query)
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

func TestDeleteQuote(t *testing.T) {
	tm, err := createTestModel(Quote{ID: 1, Text: "haskell joke"})
	defer tm.Close()
	if err != nil {
		t.Error("Got unexpected error: ", err)
	}

	q, err := tm.m.GetQuote(1)
	if err != nil {
		t.Error("Got unexpected error: ", err)
	}
	if reflect.DeepEqual(q, Quote{}) {
		t.Error("Quote should be non-zero-value: ", q)
	}

	err = tm.m.DeleteQuote(1)
	if err != nil {
		t.Error("Got unexpected error: ", err)
	}

	qu, err := tm.m.GetQuote(1)
	if err == nil {
		t.Error("got nil when error expected")
	}
	if !reflect.DeepEqual(qu, Quote{}) {
		t.Error("Quote should be zero value", q)
	}

	rawQ := rawQuote{}

	err = tm.m.db.QueryRow(
		"SELECT id, text, score, time_created, is_offensive, is_nishbot from deleted_quote where id = 1",
	).Scan(
		&rawQ.ID,
		&rawQ.Text,
		&rawQ.Score,
		&rawQ.TimeCreated,
		&rawQ.IsOffensive,
		&rawQ.IsNishbot,
	)
	if err != nil {
		t.Error("Got unexpected error: ", err)
	}
	expected := fromQuote(q)
	expected.ID = 1

	if !reflect.DeepEqual(rawQ, expected) {
		t.Error("Quote should be in deleted_quote", rawQ, expected)
	}
}

func TestVoteQuote(t *testing.T) {
	tm, err := createTestModel()
	defer tm.Close()
	if err != nil {
		t.Error("Got unexpected error: ", err)
	}

	tests := []struct {
		votes []int
		score int
	}{
		{[]int{1}, 1},
		{[]int{-1}, -1},
		{[]int{1, -1}, 0},
		{[]int{1, 1}, 2},
		{[]int{-1, -1}, -2},
	}

	for i, test := range tests {
		id := i + 1 // because 0 is a crap id
		// add new score 0 quote
		tm.m.AddQuote(Quote{ID: id, Score: 0})

		// apply votes
		for _, vote := range test.votes {
			tm.m.VoteQuote(id, vote)
		}

		// check score
		q, err := tm.m.GetQuote(id)
		if err != nil {
			t.Error("Got unexpected error: ", err)
		}

		if q.Score != test.score {
			t.Error("Got: ", q.Score, ", Expected: ", test.score)
		}
	}
}

func TestCountQuotes(t *testing.T) {
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
		search string
		count  int
	}{
		{"", 3},
		{"python", 1},
	}

	for _, test := range tests {
		count, err := tm.m.CountQuotes(test.search)
		if err != nil {
			t.Error("Got unexpected error: ", err)
		}
		if count != test.count {
			t.Error("Got: ", count, ", Expected: ", test.count)
		}
	}
}
