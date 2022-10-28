-- +goose Up
-- +goose StatementBegin
CREATE TABLE route256.financial_bot.category
(
    id      TEXT PRIMARY KEY,
    name_ru TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS route256.financial_bot.category;
-- +goose StatementEnd
