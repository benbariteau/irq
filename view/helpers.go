package view

func maxPage(totalItems, perPage int) int {
	return (totalItems - 1)/perPage + 1
}
