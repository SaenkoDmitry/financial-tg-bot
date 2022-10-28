package constants

import "github.com/pkg/errors"

const (
	WeekPeriod  = "Неделя"
	MonthPeriod = "Месяц"
	YearPeriod  = "Год"
)

const (
	ServerCurrency = "RUB"
)

const (
	Start            = "start"
	AddOperation     = "add_operation"
	SetLimitation    = "set_category_limitation"
	ShowCategoryList = "show_category_list"
	ChangeCurrency   = "change_currency"
	ShowReport       = "show_report"
)

const (
	IncorrectAmountClientMsg       = "не могу распознать введенную сумму, \n формат записи: 12345 (без пробелов и знаков препинания)"
	TransactionAddedMsg            = "Трата в категории '%s' на сумму %s %s добавлена!"
	LimitExceededMsg               = "Трата в категории '%s' на сумму %s %s добавлена, но лимит на текущий месяц превышен на %s %s !"
	SpecifyAmountMsg               = "укажите сумму расхода (%s): "
	SpecifyCategoryMsg             = "Выберите категорию:"
	SpecifyPeriodMsg               = "Выберите желаемый период:"
	SpecifyCurrencyMsg             = "Выберите валюту по умолчанию:"
	UnrecognizedCommandMsg         = "Неизвестная команда"
	SetLimitUntilDateMsg           = "Установлен лимит \nв категории '%s' \nна %s %s до даты: %s !"
	InternalServerErrorMsg         = "Внутренняя ошибка сервера"
	CannotShowCurrencyMenuMsg      = "Не могу отобразить список валют из-за внутренней ошибки :("
	HelloMsg                       = "привет, друг!"
	UndefinedCurrencyMsg           = "Бот не поддерживает выбранную вами валюту :("
	CannotChangeCurrencyMsg        = "Не могу поменять валюту :("
	CurrencyChangedSuccessfullyMsg = "Валюта успешно изменена на '%s'!"
	CannotGetRateForYouMsg         = "не могу загрузить курс из-за внутренней ошибки \xF0\x9F\x98\x94\nПопробуйте позже или выберите дефолтную валюту: %s"
	ServerProblemMsg               = "Проблемы на сервере, уже чиним \xF0\x9F\x99\x88\n\nПоказаны результаты в базовой валюте:\n\n"
)

var MissingCurrencyErr = errors.New("missing currency")
