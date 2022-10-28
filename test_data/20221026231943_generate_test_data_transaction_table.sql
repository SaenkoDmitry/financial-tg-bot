-- +goose Up
-- +goose StatementBegin
-- random category_id
INSERT INTO route256.financial_bot.user (id, currency_id)
VALUES (123548568, 'USD'),
       (1, 'RUB'),
       (12, 'USD'),
       (123, 'RUB'),
       (1234, 'EUR'),
       (12345, 'USD'),
       (123456, 'RUB'),
       (1234567, 'USD'),
       (12345678, 'EUR')
ON CONFLICT DO NOTHING;

Create or replace function random_category_id() returns text as
$$
declare
begin
    return (SELECT id FROM route256.financial_bot.category ORDER BY random() LIMIT 1);
end;
$$ language plpgsql;


-- random date
Create or replace function random_date() returns timestamp as
$$
declare
begin
    return (SELECT (ARRAY [
        '2022-05-01 23:04:05.588718',
        '2022-10-01 23:04:05.588718',
        '2022-10-05 23:04:05.588718',
        '2022-10-16 23:04:05.588718',
        '2022-10-17 01:17:04.308088'
        ])[floor(random() * 5 + 1)]);
end;
$$ language plpgsql;



-- random amount
Create or replace function random_amount() returns int as
$$
declare
begin
    return (SELECT floor(random() * 10000));
end;
$$ language plpgsql;

-- random user
Create or replace function random_user_id() returns int as
$$
declare
begin
    return (SELECT (ARRAY [
        123548568,
        1,
        12,
        123,
        1234,
        12345,
        123456,
        1234567,
        12345678
        ])[floor(random() * 9 + 1)]);
end;
$$ language plpgsql;


INSERT INTO route256.financial_bot.transaction (category_id, user_id, amount, created_at)
SELECT random_category_id() AS category_id,
       random_user_id()     AS user_id,
       random_amount()      AS amount,
       random_date()        AS created_at
FROM generate_series(1, 100000) s(i)
ON CONFLICT DO NOTHING
RETURNING (category_id, user_id, amount, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION random_category_id();
DROP FUNCTION random_date();
DROP FUNCTION random_user_id();
DROP FUNCTION random_amount();
-- +goose StatementEnd
