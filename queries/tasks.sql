-- name: CreateTaskWithID :one
INSERT INTO tasks (id, event)
VALUES (?, ?)
RETURNING *;

-- name: CreateTask :one
INSERT INTO tasks (event)
VALUES (?)
RETURNING *;

-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY id ASC;
