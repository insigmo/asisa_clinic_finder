-- +goose Up
CREATE TABLE IF NOT EXISTS user (
    id bigint primary key,
    username text,
    name text,
    lastname text,
    is_bot bool,
    city text,
    language_code text
);

-- +goose Down
DROP TABLE IF EXISTS user;
