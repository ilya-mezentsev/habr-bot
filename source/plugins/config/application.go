package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
)

type Configs struct {
	DBPath               string
	Categories           []string `json:"categories"`
	CategoriesFilters    []string `json:"categories_filters"`
	ArticlesResource     string   `json:"articles_resource"`
	ArticleLinkClassName string   `json:"article_link_class_name"`
	ArticlesFilter       string   `json:"articles_filter"`
	Mode                 Mode
}

func GetAll() (configs Configs, err error) {
	if !flag.Parsed() {
		flag.Parse()
	}

	data, err := ioutil.ReadFile(*configsPath)
	if err != nil {
		return Configs{}, err
	}

	err = json.Unmarshal(data, &configs)
	if err != nil {
		return Configs{}, err
	}

	configs.DBPath, err = GetDBPath()
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

func GetTelegramBotToken() (string, error) {
	token := os.Getenv("TG_BOT_TOKEN")
	if token == "" {
		return "", noTelegramBotToken
	}

	return token, nil
}

func GetMode() Mode {
	return mode
}
