package commands

import "strings"

const (
	GetCategories = "/get_categories"
	GetArticles   = "/get_articles"
	ParseAll      = "/parse_all"
	ParseArticles = "/parse_articles"
)

func All() string {
	return strings.Join([]string{
		GetCategories,
		GetArticles,
		ParseArticles,
		ParseAll,
	}, ", \n")
}
