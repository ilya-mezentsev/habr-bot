package main

import (
	cliController "controllers/cli"
	"controllers/telegram"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"interfaces"
	"mock"
	"models"
	"net/http"
	"net/url"
	"os"
	"plugins/config"
	"plugins/logger"
	cliPresenter "presenters/cli"
	tgPresenter "presenters/telegram"
	articlesRepository "repositories/articles"
	articlesService "services/articles"
	articlesParser "services/articles_parser"
)

var (
	configs    config.Configs
	controller interfaces.Controller
)

func init() {
	var err error
	configs, err = config.GetAll()
	handleError(err)

	db, err := sqlx.Open("sqlite3", configs.DBPath)
	handleError(err)
	mock.CreateTableIfNotExists(db)

	controller = getController(
		articlesService.New(
			articlesRepository.New(db),
			articlesParser.New(
				configs.ArticlesResource,
				configs.ArticleLinkClassName,
				configs.ArticlesFilter,
			),
			configs.Categories,
		),
	)
}

func getController(service interfaces.ArticlesService) interfaces.Controller {
	switch {
	case configs.Mode.IsCLI():
		return cliController.New(
			service,
			cliPresenter.New(),
		)
	case configs.Mode.IsTelegram():
		return getTelegramController(service)
	default:
		logger.Error("Unknown mode")
		os.Exit(1)
		return nil
	}
}

func getTelegramController(service interfaces.ArticlesService) interfaces.Controller {
	tgToken, err := config.GetTelegramBotToken()
	handleError(err)

	httpProxyURL, err := config.GetHttpProxyURL()
	handleError(err)

	proxyUrl, err := url.Parse(fmt.Sprintf("http://%s", httpProxyURL))
	handleError(err)
	http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

	bot, err := tg.NewBotAPI(tgToken)
	handleError(err)

	return telegram.New(
		service,
		tgPresenter.New(bot),
		bot,
	)
}

func handleError(err error) {
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func main() {
	logger.Info("Start application in mode: " + configs.Mode.Value())
	processing := models.ProcessingChannels{
		Error: make(chan error),
	}

	go controller.Run(processing)

	for {
		select {
		case err := <-processing.Error:
			logger.Error(err)
			return
		}
	}
}
