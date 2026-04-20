-- +goose Up
CREATE TABLE IF NOT EXISTS medical_direction (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    reference_name TEXT NOT NULL,
    name           TEXT    NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_directions ON medical_direction (reference_name);

-- +goose Down
DROP TABLE IF EXISTS medical_direction;
DROP INDEX IF EXISTS idx_directions;
