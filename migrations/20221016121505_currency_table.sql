-- +goose Up
-- +goose StatementBegin
CREATE TABLE route256.financial_bot.currency
(
    id      TEXT PRIMARY KEY UNIQUE,
    name_ru TEXT NOT NULL,
    symbol  TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS route256.financial_bot.currency;
-- +goose StatementEnd
