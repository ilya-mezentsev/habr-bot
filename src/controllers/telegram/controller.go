package telegram

import (
	"controllers/commands"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"interfaces"
	"models"
	"services/errors"
)

var currentCategory = ""

type Controller struct {
	service   interfaces.ArticlesService
	presenter interfaces.TelegramPresenter
	bot       *tg.BotAPI
}

func New(
	service interfaces.ArticlesService,
	presenter interfaces.TelegramPresenter,
	bot *tg.BotAPI,
) Controller {
	return Controller{service, presenter, bot}
}

func (c Controller) Run(processing models.ProcessingChannels) {
	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates, err := c.bot.GetUpdatesChan(u)
	if err != nil {
		processing.Error <- err
		return
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			err := c.processChoseCategory(update)
			if err != nil {
				processing.Error <- err
				return
			}
		}
		if update.Message != nil {
			err := c.processMessage(update)
			if err != nil {
				processing.Error <- err
				return
			}
		}
	}
}

func (c Controller) processChoseCategory(update tg.Update) error {
	currentCategory = update.CallbackQuery.Data
	_, err := c.bot.AnswerCallbackQuery(tg.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
	if err != nil {
		return err
	}

	err = c.presenter.Info(tg.NewMessage(
		update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Chosen cateogry: %s", update.CallbackQuery.Data),
	))
	if err != nil {
		return err
	}

	return nil
}

func (c Controller) processMessage(update tg.Update) error {
	switch update.Message.Text {
	case commands.ParseAll:
		return c.parseAll(update)
	case commands.GetCategories:
		return c.showCategories(update)
	case commands.GetArticles:
		return c.showArticles(update)
	case commands.ParseArticles:
		return c.parseArticles(update)
	default:
		return c.sayAboutUnknownCommandAndShowCategories(update)
	}
}

func (c Controller) parseAll(update tg.Update) error {
	err := c.service.ParseAll()
	if err != nil {
		return err
	}

	return c.sayAllCategoriesIsParsed(update)
}

func (c Controller) sayAboutUnknownCommandAndShowCategories(update tg.Update) error {
	err := c.presenter.Info(tg.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("Unknown command. Available: \n%s", commands.All()),
	))
	if err != nil {
		return err
	}

	return c.showCategories(update)
}

func (c Controller) showCategories(update tg.Update) error {
	c.presenter.SetMessageConfig(tg.NewMessage(update.Message.Chat.ID, "Choose category"))
	err := c.presenter.ShowCategories(c.service.GetCategories())
	if err != nil {
		return err
	}

	return nil
}

func (c Controller) showArticles(update tg.Update) error {
	if currentCategory == "" {
		return c.sayCategoryIsNotSet(update)
	}

	articles, err := c.service.GetArticles(currentCategory)
	switch err {
	case nil:
		break
	case errors.NoArticlesForCategory:
		return c.sayNoArticlesForCategory(update)
	default:
		return err
	}

	c.presenter.SetMessageConfig(tg.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("Articles for %s category", currentCategory),
	))
	return c.presenter.ShowArticles(articles)
}

func (c Controller) parseArticles(update tg.Update) error {
	if currentCategory == "" {
		return c.sayCategoryIsNotSet(update)
	}

	err := c.service.ParseCategory(currentCategory)
	if err != nil {
		return err
	}

	return c.sayArticlesIsParsed(update)
}

func (c Controller) sayCategoryIsNotSet(update tg.Update) error {
	return c.presenter.Info(tg.NewMessage(
		update.Message.Chat.ID,
		"Category is not set",
	))
}

func (c Controller) sayNoArticlesForCategory(update tg.Update) error {
	return c.presenter.Info(tg.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("No articles for category: %s", currentCategory),
	))
}

func (c Controller) sayArticlesIsParsed(update tg.Update) error {
	return c.presenter.Info(tg.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("Articles is parsed for category: %s", currentCategory),
	))
}

func (c Controller) sayAllCategoriesIsParsed(update tg.Update) error {
	return c.presenter.Info(tg.NewMessage(
		update.Message.Chat.ID,
		"All categories is parsed",
	))
}
