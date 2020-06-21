package telegram

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"interfaces"
	"models"
)

type Presenter struct {
	bot           *tg.BotAPI
	messageConfig tg.MessageConfig
}

func New(bot *tg.BotAPI) interfaces.TelegramPresenter {
	return &Presenter{bot, tg.MessageConfig{}}
}

func (p *Presenter) SetMessageConfig(config tg.MessageConfig) {
	p.messageConfig = config
}

func (p Presenter) ShowAsButtons(categories []string) error {
	if p.messageConfig == (tg.MessageConfig{}) {
		return noMessageConfig
	}

	p.messageConfig.ReplyMarkup = tg.NewInlineKeyboardMarkup(p.getCategoriesMarkup(categories)...)
	_, err := p.bot.Send(p.messageConfig)

	return err
}

func (p Presenter) getCategoriesMarkup(categories []string) [][]tg.InlineKeyboardButton {
	var inlineButtons [][]tg.InlineKeyboardButton
	for i := 0; i < len(categories); i++ {
		buttons := tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(categories[i], categories[i]),
		)
		if i < len(categories)-1 {
			buttons = append(buttons, tg.NewInlineKeyboardButtonData(categories[i+1], categories[i+1]))
			i++
		}

		inlineButtons = append(inlineButtons, buttons)
	}

	return inlineButtons
}

func (p Presenter) ShowArticles(articles []models.Article) error {
	if p.messageConfig == (tg.MessageConfig{}) {
		return noMessageConfig
	}

	p.messageConfig.ReplyMarkup = tg.NewInlineKeyboardMarkup(p.getArticlesMarkup(articles)...)
	_, err := p.bot.Send(p.messageConfig)

	return err
}

func (p Presenter) getArticlesMarkup(articles []models.Article) [][]tg.InlineKeyboardButton {
	var buttons [][]tg.InlineKeyboardButton
	for _, article := range articles {
		buttons = append(buttons, tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonURL(article.Title, article.Link),
		))
	}

	return buttons
}

func (p Presenter) Info(message tg.MessageConfig) error {
	_, err := p.bot.Send(message)
	return err
}
