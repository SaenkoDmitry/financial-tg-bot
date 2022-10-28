-- +goose Up
-- +goose StatementBegin
INSERT INTO route256.financial_bot.category (id, name_ru)
VALUES ('FASTFOOD', '🍔 Фаст-фуд'),
       ('RESTAURANTS', '🍷 Рестораны'),
       ('SUPERMARKETS', '🏪 Супермаркеты'),
       ('CLOTHES', '👔 Одежда'),
       ('EDUCATION', '🎓 Образование'),
       ('TRANSPORT', '🚕 Транспорт'),
       ('MEDICINE', '💊 Медицина'),
       ('BEAUTY', '💅 Красота'),
       ('ENTERTAINMENT', '🎡 Развлечения'),
       ('UNSCHEDULED', '🕐 Незапланированное'),
       ('OTHERS', '💸 Другое');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE route256.financial_bot.category;
-- +goose StatementEnd
