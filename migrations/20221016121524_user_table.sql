-- +goose Up
-- +goose StatementBegin
CREATE TABLE route256.financial_bot.user
(
    id          BIGINT PRIMARY KEY,
    currency_id TEXT NOT NULL DEFAULT 'RUB' REFERENCES route256.financial_bot.currency (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS route256.financial_bot.user;
-- +goose StatementEnd
