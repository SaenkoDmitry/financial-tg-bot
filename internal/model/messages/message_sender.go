package messages

import (
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
	SendMessageWithMarkup(text string, markup [][]model.MarkupData, userID int64) error
}
