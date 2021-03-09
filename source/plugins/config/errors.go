package config

import "errors"

var (
	noDBFile               = errors.New("no DB file")
	noCategories           = errors.New("no categories")
	noFilters              = errors.New("no categories filters")
	noArticlesResource     = errors.New("no articles resource")
	noTelegramBotToken     = errors.New("no telegram token")
	noArticleLinkClassName = errors.New("no article link class name")
)
