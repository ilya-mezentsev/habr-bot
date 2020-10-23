package commands

import "strings"

const (
	GetCategories = "/get_categories"
	GetFilters    = "/get_filters"
	GetArticles   = "/get_articles"
	ParseAll      = "/parse_all"
	ParseArticles = "/parse_articles"
	Clean         = "/clean"
)

func All() string {
	return strings.Join([]string{
		GetCategories,
		GetFilters,
		GetArticles,
		ParseArticles,
		ParseAll,
		Clean,
	}, "\n")
}
