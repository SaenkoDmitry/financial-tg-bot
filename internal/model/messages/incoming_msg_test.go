package messages

import (
	"github.com/golang/mock/gomock"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/constants"
	mocks "gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/mocks/messages"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CurrencyGetterMock struct {
}

func (c *CurrencyGetterMock) Currencies() []string {
	return []string{"RUB", "USD", "EUR", "CNY"}
}

func TestOnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	currencySetter := &CurrencyGetterMock{}
	sender := mocks.NewMockMessageSender(ctrl)
	userCurrencyRepo, _ := repository.NewUserCurrencyRepository(currencySetter)
	model := New(sender, userCurrencyRepo)

	sender.EXPECT().SendMessage(constants.HelloMsg, int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/start",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func TestOnStartCommand_ShouldAnswerWithUnexpectedMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	currencySetter := &CurrencyGetterMock{}
	sender := mocks.NewMockMessageSender(ctrl)
	userCurrencyRepo, _ := repository.NewUserCurrencyRepository(currencySetter)
	model := New(sender, userCurrencyRepo)

	sender.EXPECT().SendMessage(constants.UnrecognizedCommandMsg, int64(123))

	err := model.IncomingMessage(Message{
		Text:   "what?",
		UserID: 123,
	})

	assert.NoError(t, err)
}
