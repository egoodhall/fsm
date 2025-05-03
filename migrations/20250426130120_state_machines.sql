-- +goose Up
-- +goose StatementBegin
CREATE TABLE state_machines (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at INTEGER NOT NULL DEFAULT(unixepoch('subsec') * 1000)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE state_machines;
-- +goose StatementEnd
