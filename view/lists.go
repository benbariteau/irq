package view

var (
	All    = QuotesBase("All", []string{"id ASC"})
	Latest = QuotesBase("Latest", []string{"id DESC"})
	Search = QuotesBase("Search", []string{"score DESC"})
	Top    = QuotesBase("Top", []string{"score DESC"})
)
