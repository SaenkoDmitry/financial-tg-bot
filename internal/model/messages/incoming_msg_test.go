package messages

import (
	"github.com/golang/mock/gomock"
	mocks "gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/mocks/messages"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender)

	sender.EXPECT().SendMessage("hello", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/start",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func TestOnStartCommand_ShouldAnswerWithUnexpectedMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender)

	sender.EXPECT().SendMessage("unexpected command", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "what?",
		UserID: 123,
	})

	assert.NoError(t, err)
}
