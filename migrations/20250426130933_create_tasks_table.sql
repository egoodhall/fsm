-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
    id TEXT NOT NULL PRIMARY KEY,
    event BLOB NOT NULL,
    fsm_id INTEGER NOT NULL REFERENCES state_machines(id),
    created_at INTEGER NOT NULL DEFAULT(unixepoch('subsec') * 1000)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;
-- +goose StatementEnd
