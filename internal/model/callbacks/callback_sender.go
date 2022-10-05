package callbacks

import "gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model"

type CallbackSender interface {
	SendMessage(text string, userID int64) error
	SendMessageWithMarkup(text string, markup [][]model.MarkupData, userID int64) error
	SendEditMessage(text string, userID int64, messageID int) error
	SendEditMessageWithMarkupAndText(text string, markup [][]model.MarkupData, userID int64, messageID int) error
}
