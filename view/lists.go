package view

var (
	All    = QuotesBase("All", []string{"quote.id ASC"})
	Latest = QuotesBase("Latest", []string{"quote.id DESC"})
	Search = QuotesBase("Search", []string{"quote.score DESC"})
	Top    = QuotesBase("Top", []string{"quote.score DESC"})
)
