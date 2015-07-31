package model

import (
	"time"
)

func searchWhereClause(search string) string {
	if search == "" {
		return ""
	}
	return "WHERE text LIKE '%" + search + "%'"
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
