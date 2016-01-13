package model

import (
	"strings"
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
		Tags:        strings.Split(strings.ToLower(rawQ.Tags), ","),
	}
}

func fromQuote(quote Quote) rawQuote {
	return rawQuote{
		ID:          quote.ID,
		Text:        quote.Text,
		Score:       quote.Score,
		TimeCreated: quote.TimeCreated,
	}
}
