-- +goose Up
-- +goose StatementBegin
CREATE TABLE state_transitions (
    id INTEGER PRIMARY KEY,
    attempt INTEGER NOT NULL,
    task_id INTEGER NOT NULL REFERENCES tasks(id),
    from_state TEXT NOT NULL,
    to_state TEXT NOT NULL,
    data BLOB DEFAULT NULL,
    created_at INTEGER NOT NULL DEFAULT(unixepoch('subsec') * 1000)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE state_transitions;
-- +goose StatementEnd
