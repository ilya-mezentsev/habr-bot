package config

import (
	"flag"
	"os"
	"strings"
)

type Configs struct {
	DBPath               string
	Categories           []string
	Filters              []string
	ArticlesResource     string
	ArticleLinkClassName string
	ArticlesFilter       string
	Mode                 Mode
}

func GetAll() (configs Configs, err error) {
	configs.DBPath, err = GetDBPath()
	if err != nil {
		return Configs{}, err
	}

	configs.Categories, err = GetCategories()
	if err != nil {
		return Configs{}, err
	}

	configs.Filters, err = GetFilters()
	if err != nil {
		return Configs{}, err
	}

	configs.ArticlesResource, err = GetArticleResource()
	if err != nil {
		return Configs{}, err
	}

	configs.ArticleLinkClassName, err = GetArticleLinkClassName()
	if err != nil {
		return Configs{}, err
	}

	configs.Mode = GetMode()
	return configs, nil
}

func GetDBPath() (string, error) {
	path := os.Getenv("DB_FILE")
	if path == "" {
		return "", noDBFile
	}

	return path, nil
}

func GetCategories() ([]string, error) {
	categories := os.Getenv("CATEGORIES")
	if categories == "" {
		return nil, noCategories
	}

	return strings.Split(categories, ","), nil
}

func GetFilters() ([]string, error) {
	filters := os.Getenv("CATEGORIES_FILTERS")
	if filters == "" {
		return nil, noFilters
	}

	return strings.Split(filters, ","), nil
}

func GetArticleResource() (string, error) {
	resource := os.Getenv("ARTICLES_RESOURCE")
	if resource == "" {
		return "", noArticlesResource
	}

	return resource, nil
}

func GetArticleLinkClassName() (string, error) {
	resource := os.Getenv("ARTICLE_LINK_CLASS_NAME")
	if resource == "" {
		return "", noArticleLinkClassName
	}

	return resource, nil
}

func GetTelegramBotToken() (string, error) {
	token := os.Getenv("TG_BOT_TOKEN")
	if token == "" {
		return "", noTelegramBotToken
	}

	return token, nil
}

func GetMode() Mode {
	flag.Parse()
	return mode
}
