package main

import (
	cliController "controllers/cli"
	"controllers/telegram"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"interfaces"
	"mock"
	"models"
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
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	db, err := sqlx.Open("sqlite3", configs.DBPath)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	mock.CreateTableIfNotExists(db)

	bot, err := tg.NewBotAPI(configs.TgToken)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

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
		bot,
	)
}

func getController(service interfaces.ArticlesService, bot *tg.BotAPI) interfaces.Controller {
	switch {
	case configs.Mode.IsCLI():
		return cliController.New(
			service,
			cliPresenter.New(),
		)
	case configs.Mode.IsTelegram():
		return telegram.New(
			service,
			tgPresenter.New(bot),
			bot,
		)
	default:
		logger.Error("Unknown mode")
		os.Exit(1)
		return nil
	}
}

func main() {
	logger.Info("Start application in mode: " + configs.Mode.Value())
	processing := models.ProcessingChannels{
		Done:  make(chan bool),
		Error: make(chan error),
	}

	go controller.Run(processing)

	for {
		select {
		case err := <-processing.Error:
			logger.Error(err)
			return
		case <-processing.Done:
			logger.Info("Processing done")
			return
		}
	}
}
