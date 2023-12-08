-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
TRUNCATE TABLE users;
TRUNCATE TABLE  user_balance;
TRUNCATE TABLE  user_orders;
TRUNCATE TABLE billing;