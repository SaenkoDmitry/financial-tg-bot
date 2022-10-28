-- +goose Up
-- +goose StatementBegin
CREATE TABLE route256.financial_bot.rate
(
    currency_id TEXT    NOT NULL REFERENCES route256.financial_bot.currency (id),
    multiplier  DECIMAL NOT NULL,
    on_date     DATE    NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS route256.financial_bot.rate;
-- +goose StatementEnd
