package interfaces

import "habr-bot/source/models"

type (
	ArticlesRepository interface {
		Save(articles chan models.Article, trySave chan bool, processing models.ProcessingChannels)
		GetByCategory(category string) ([]models.Article, error)
		ClearArticles() error
		ClearCategoryArticles(category string) error
	}
)
