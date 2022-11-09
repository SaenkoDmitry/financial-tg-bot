package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/config"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model/callbacks"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model/messages"
	"go.uber.org/zap"
)

type Client struct {
	client *tgbotapi.BotAPI
}

type TokenGetter interface {
	Token() string
}

func New(cfg *config.Service) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(cfg.Token())
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	if _, err = client.Request(initialCommands); err != nil {
		return nil, errors.Wrap(err, "cannot set methods")
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) SendMessage(text string, userID int64) error {
	_, err := c.client.Send(tgbotapi.NewMessage(userID, text))
	if err != nil {
		return errors.Wrap(err, "cannot execute SendMessage")
	}
	return nil
}

func (c *Client) SendMessageWithMarkup(text string, markup [][]model.MarkupData, userID int64) error {
	msg := tgbotapi.NewMessage(userID, text)
	msg.ReplyMarkup = buildReplyMarkup(markup)
	_, err := c.client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "cannot execute SendMessageWithMarkup")
	}
	return nil
}

func buildReplyMarkup(markup [][]model.MarkupData) tgbotapi.InlineKeyboardMarkup {
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, len(markup))
	for i := range markup {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(lo.Map(markup[i], mapMarkup)...))
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func mapMarkup(t model.MarkupData, _ int) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(t.Text, t.Data)
}

func (c *Client) SendEditMessageWithMarkupAndText(text string, markup [][]model.MarkupData,
	userID int64, messageID int) error {
	replyMarkup := buildReplyMarkup(markup)
	_, err := c.client.Send(tgbotapi.NewEditMessageTextAndMarkup(userID, messageID, text, replyMarkup))
	if err != nil {
		return errors.Wrap(err, "cannot execute SendEditMessageWithMarkupAndText")
	}
	return nil
}

func (c *Client) SendEditMessage(text string, userID int64, messageID int) error {
	_, err := c.client.Send(tgbotapi.NewEditMessageText(userID, messageID, text))
	if err != nil {
		return errors.Wrap(err, "cannot execute SendEditMessage")
	}
	return nil
}

func (c *Client) ListenUpdates(ctx context.Context, msgModel *messages.Model, callbackModel *callbacks.Model) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)
	logger.Info("start listening telegram server for new messages")

	for update := range updates {
		if update.CallbackQuery != nil {
			err := callbackModel.HandleIncomingCallback(ctx, update.CallbackQuery)
			if err != nil {
				logger.Error("error occurred while processing callback", zap.Error(err))
				continue
			}
		}
		if update.Message != nil {
			err := msgModel.IncomingMessage(ctx, messages.Message{
				Text:   update.Message.Text,
				UserID: update.Message.From.ID,
			})
			if err != nil {
				logger.Error("error occurred while processing message", zap.Error(err))
				continue
			}
		}
	}
}

var initialCommands = tgbotapi.NewSetMyCommands(
	tgbotapi.BotCommand{
		Command:     constants.AddOperation,
		Description: "добавить новую трату",
	},
	tgbotapi.BotCommand{
		Command:     constants.ShowCategoryList,
		Description: "показать список категорий",
	},
	tgbotapi.BotCommand{
		Command:     constants.SetLimitation,
		Description: "установить лимит трат (месяц)",
	},
	tgbotapi.BotCommand{
		Command:     constants.ChangeCurrency,
		Description: "сменить валюту",
	},
	tgbotapi.BotCommand{
		Command:     constants.ShowReport,
		Description: "показать отчет о тратах за период",
	},
)
