-- +goose Up
-- +goose StatementBegin
CREATE TABLE route256.financial_bot.transaction
(
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    category_id TEXT      NOT NULL REFERENCES route256.financial_bot.category (id),
    user_id     INT       NOT NULL REFERENCES route256.financial_bot.user (id),
    amount      DECIMAL   NOT NULL,
    created_at  TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS route256.financial_bot.transaction;
-- +goose StatementEnd
