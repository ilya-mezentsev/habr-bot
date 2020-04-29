package interfaces

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"models"
)

type (
	ArticlesPresenter interface {
		ShowCategories(categories []string) error
		ShowArticles(articles []models.Article) error
	}

	TelegramPresenter interface {
		ArticlesPresenter
		SetMessageConfig(message tg.MessageConfig)
		Info(message tg.MessageConfig) error
	}
)
