package telegram

import (
	"controllers/commands"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"interfaces"
	"models"
	"services/errors"
)

var (
	currentCategory   = ""
	currentFilter     = ""
	processingCommand = ""
)

type Controller struct {
	service       interfaces.ArticlesService
	presenter     interfaces.TelegramPresenter
	bot           *tg.BotAPI
	buildCategory interfaces.BuildFormattedCategory
}

func New(
	service interfaces.ArticlesService,
	presenter interfaces.TelegramPresenter,
	bot *tg.BotAPI,
	buildCategory interfaces.BuildFormattedCategory,
) Controller {
	return Controller{service, presenter, bot, buildCategory}
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
			err := c.processCallback(update)
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

func (c Controller) processCallback(update tg.Update) error {
	switch processingCommand {
	case commands.GetCategories:
		return c.processChoseCategory(update)
	case commands.GetFilters:
		return c.processChoseFilter(update)
	default:
		return nil
	}
}

func (c Controller) processChoseCategory(update tg.Update) error {
	currentCategory = update.CallbackQuery.Data
	_, err := c.bot.AnswerCallbackQuery(tg.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
	if err != nil {
		return err
	}

	return c.presenter.Info(tg.NewMessage(
		update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Chosen cateogry: %s", currentCategory),
	))
}

func (c Controller) processChoseFilter(update tg.Update) error {
	currentFilter = update.CallbackQuery.Data
	_, err := c.bot.AnswerCallbackQuery(tg.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
	if err != nil {
		return err
	}

	return c.presenter.Info(tg.NewMessage(
		update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Chosen filter: %s", currentFilter),
	))
}

func (c Controller) processMessage(update tg.Update) error {
	processingCommand = update.Message.Text

	switch processingCommand {
	case commands.ParseAll:
		return c.parseAll(update)
	case commands.GetCategories:
		return c.showCategories(update)
	case commands.GetFilters:
		return c.showFilters(update)
	case commands.GetArticles:
		return c.showArticles(update)
	case commands.ParseArticles:
		return c.parseArticles(update)
	case commands.Clean:
		return c.clean(update)
	default:
		processingCommand = ""
		return c.sayAboutUnknownCommand(update)
	}
}

func (c Controller) parseAll(update tg.Update) error {
	err := c.service.ParseAll()
	if err != nil {
		return err
	}

	return c.sayAllCategoriesIsParsed(update)
}

func (c Controller) sayAboutUnknownCommand(update tg.Update) error {
	return c.presenter.Info(tg.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("Unknown command. Available: \n%s", commands.All()),
	))
}

func (c Controller) showCategories(update tg.Update) error {
	c.presenter.SetMessageConfig(tg.NewMessage(update.Message.Chat.ID, "Choose category"))

	return c.presenter.ShowAsButtons(c.service.GetCategories())
}

func (c Controller) showFilters(update tg.Update) error {
	c.presenter.SetMessageConfig(tg.NewMessage(update.Message.Chat.ID, "Choose category"))

	return c.presenter.ShowAsButtons(c.service.GetFilters())
}

func (c Controller) showArticles(update tg.Update) error {
	if currentCategory == "" {
		return c.sayCategoryIsNotSet(update)
	}

	category := c.buildCategory(currentCategory, currentFilter)
	articles, err := c.service.GetArticles(category)
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
		fmt.Sprintf("Articles for %s category", category),
	))
	return c.presenter.ShowArticles(articles)
}

func (c Controller) parseArticles(update tg.Update) error {
	if currentCategory == "" {
		return c.sayCategoryIsNotSet(update)
	}

	category := c.buildCategory(currentCategory, currentFilter)
	err := c.service.ParseCategory(category)
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

func (c Controller) clean(update tg.Update) error {
	currentCategory = ""
	currentFilter = ""

	return c.sayCategoryAndFilterAreCleaned(update)
}

func (c Controller) sayCategoryAndFilterAreCleaned(update tg.Update) error {
	return c.presenter.Info(tg.NewMessage(
		update.Message.Chat.ID,
		"Category and filter are cleaned",
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
