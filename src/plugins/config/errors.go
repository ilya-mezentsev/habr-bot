package config

import "errors"

var (
	noDBFile               = errors.New("no DB file")
	noCategories           = errors.New("no categories")
	noArticlesResource     = errors.New("no articles resource")
	noTelegramBotToken     = errors.New("no telegram token")
	noArticleLinkClassName = errors.New("no article link class name")
)
