package messages

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	messagesMocks "gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/mocks/messages"

	"github.com/stretchr/testify/assert"
)

func TestOnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	ctx := context.Background()
	sender := messagesMocks.NewMockMessageSender(ctrl)
	userRepoMock := messagesMocks.NewMockUserStore(ctrl)
	categoryRepoMock := messagesMocks.NewMockCategoryStore(ctrl)
	model := New(sender, userRepoMock, categoryRepoMock)

	userRepoMock.EXPECT().SetUserCurrency(gomock.Any(), int64(123), "RUB").Times(1)
	sender.EXPECT().SendMessage(constants.HelloMsg, int64(123))

	err := model.IncomingMessage(ctx, Message{
		Text:   "/start",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func TestOnStartCommand_ShouldAnswerWithUnexpectedMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	ctx := context.Background()
	sender := messagesMocks.NewMockMessageSender(ctrl)
	userRepoMock := messagesMocks.NewMockUserStore(ctrl)
	categoryRepoMock := messagesMocks.NewMockCategoryStore(ctrl)
	model := New(sender, userRepoMock, categoryRepoMock)

	sender.EXPECT().SendMessage(constants.UnrecognizedCommandMsg, int64(123))

	err := model.IncomingMessage(ctx, Message{
		Text:   "what?",
		UserID: 123,
	})

	assert.NoError(t, err)
}
