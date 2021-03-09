package main

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	cliController "habr-bot/source/controllers/cli"
	"habr-bot/source/controllers/telegram"
	"habr-bot/source/interfaces"
	"habr-bot/source/mock"
	"habr-bot/source/models"
	"habr-bot/source/plugins/config"
	"habr-bot/source/plugins/logger"
	cliPresenter "habr-bot/source/presenters/cli"
	tgPresenter "habr-bot/source/presenters/telegram"
	articlesRepository "habr-bot/source/repositories/articles"
	articlesService "habr-bot/source/services/articles"
	articlesParser "habr-bot/source/services/articles_parser"
	"habr-bot/source/services/category_format"
	"os"
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
			),
			configs.Categories,
			configs.Filters,
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

	bot, err := tg.NewBotAPI(tgToken)
	handleError(err)

	_, err = bot.RemoveWebhook()
	handleError(err)

	return telegram.New(
		service,
		tgPresenter.New(bot),
		bot,
		category_format.GetFormattedCategory,
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
