-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA financial_bot;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS financial_bot;
-- +goose StatementEnd
