-- +goose Up
-- +goose StatementBegin
ALTER TABLE route256.financial_bot.rate
    ADD CONSTRAINT unique_currency_id_on_date UNIQUE (currency_id, on_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE route256.financial_bot.rate
    DROP CONSTRAINT unique_currency_id_on_date;
-- +goose StatementEnd
