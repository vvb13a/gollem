-- +goose Up
-- Ground-zero schema. Later changes go in new numbered migration files.

CREATE TABLE items (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    notes TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX idx_items_created_at ON items(created_at);

-- +goose Down
DROP TABLE items;
