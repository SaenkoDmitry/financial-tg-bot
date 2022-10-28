-- +goose Up
-- +goose StatementBegin
CREATE TABLE route256.financial_bot.limitation
(
    id           INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    category_id  TEXT    NOT NULL REFERENCES route256.financial_bot.category (id),
    user_id      INT     NOT NULL REFERENCES route256.financial_bot.user (id),
    upper_border DECIMAL NOT NULL,
    until_date   DATE    NOT NULL
);

ALTER TABLE route256.financial_bot.limitation
    ADD CONSTRAINT unique_category_id_user_id UNIQUE (category_id, user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE route256.financial_bot.limitation
    DROP CONSTRAINT unique_category_id_user_id;
DROP TABLE IF EXISTS route256.financial_bot.limitation;
-- +goose StatementEnd
