-- +goose Up
CREATE TABLE IF NOT EXISTS city_postal_codes (
     id          INTEGER PRIMARY KEY AUTOINCREMENT,
     city        TEXT    NOT NULL,
     postal_code INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_city ON city_postal_codes (city);

-- +goose Down
DROP TABLE IF EXISTS city_postal_codes;
DROP INDEX IF EXISTS idx_city;