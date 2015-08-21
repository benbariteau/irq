package model

import (
	"time"
)

func toQuote(rawQ rawQuote) Quote {
	var timeCreated time.Time
	switch rawQ.TimeCreated.(type) {
	case time.Time:
		timeCreated = rawQ.TimeCreated.(time.Time)
	case int64:
		timeCreated = time.Unix(rawQ.TimeCreated.(int64), 0)
	}

	return Quote{
		ID:          rawQ.ID,
		Text:        rawQ.Text,
		Score:       rawQ.Score,
		TimeCreated: timeCreated,
		IsOffensive: rawQ.IsOffensive != 0,
		IsNishbot:   rawQ.IsNishbot != 0,
	}
}

func fromQuote(quote Quote) rawQuote {
	return rawQuote{
		ID:          quote.ID,
		Text:        quote.Text,
		Score:       quote.Score,
		TimeCreated: quote.TimeCreated,
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
