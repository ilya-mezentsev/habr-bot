package interfaces

import "models"

type (
	ArticlesService interface {
		ParseAll() error
		ParseCategory(category string) error
		GetCategories() []string
		GetArticles(category string) ([]models.Article, error)
	}

	ArticlesParserService interface {
		ParseCategory(category string, articles chan<- models.Article, processing models.ProcessingChannels)
	}
)
