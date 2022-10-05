package constants

const (
	WeekPeriod  = "Неделя"
	MonthPeriod = "Месяц"
	YearPeriod  = "Год"
)

const (
	Start            = "start"
	AddOperation     = "add_operation"
	ShowCategoryList = "show_category_list"
	ShowReport       = "show_report"
)

var CategoryList = []string{
	FastFood, Restaurants, Supermarkets, Clothes, Education, Transport, Medicine, Beauty, Entertainment, Unscheduled, Others,
}

const (
	FastFood      = "\xF0\x9F\x8D\x94 Фаст-фуд"
	Restaurants   = "\xF0\x9F\x8D\xB7 Рестораны"
	Supermarkets  = "\xF0\x9F\x8F\xAA Супермаркеты"
	Clothes       = "\xF0\x9F\x91\x95 Одежда"
	Education     = "\xF0\x9F\x8E\x93 Образование"
	Transport     = "\xF0\x9F\x9A\x95 Транспорт"
	Medicine      = "\xF0\x9F\x92\x8A Медицина"
	Beauty        = "\xF0\x9F\x92\x85 Красота"
	Entertainment = "\xF0\x9F\x8E\xAA Развлечения"
	Unscheduled   = "\xF0\x9F\x95\x90 Незапланированное"
	Others        = "\xF0\x9F\x92\xB8 Другое"
)

const (
	IncorrectAmountClientMsg = "не могу распознать введенную сумму, \n формат записи: 12345 (без пробелов и знаков препинания)"
	TransactionAddedMsg      = "Трата в категории '%s' на сумму %s руб. добавлена!"
	SpecifyAmountMsg         = "укажите сумму расхода (руб.): "
	SpecifyCategoryMsg       = "Выберите категорию:"
	SpecifyPeriodMsg         = "Выберите желаемый период:"
	UnrecognizedCommandMsg   = "Неизвестная команда"
	HelloMsg                 = "привет, друг!"
)
