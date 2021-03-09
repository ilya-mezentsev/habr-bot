package interfaces

import "habr-bot/source/models"

type (
	ArticlesService interface {
		ParseAll() error
		ParseCategory(category string) error
		GetCategories() []string
		GetFilters() []string
		GetArticles(category string) ([]models.Article, error)
	}

	ArticlesParserService interface {
		ParseCategory(category string, articles chan<- models.Article, processing models.ProcessingChannels)
	}

	BuildFormattedCategory func(name, filter string) string
)
