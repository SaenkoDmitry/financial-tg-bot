-- +goose Up
-- +goose StatementBegin
INSERT INTO route256.financial_bot.currency (id, name_ru, symbol)
VALUES ('RUB', 'руб.', '₽'),
       ('USD', 'долл.', '$'),
       ('EUR', 'евро', '€'),
       ('CNY', 'юан.', '¥');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE route256.financial_bot.currency;
-- +goose StatementEnd
