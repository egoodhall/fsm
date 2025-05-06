-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
    id INTEGER NOT NULL PRIMARY KEY,
    data BLOB NOT NULL,
    created_at INTEGER NOT NULL DEFAULT(unixepoch('subsec') * 1000)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;
-- +goose StatementEnd
