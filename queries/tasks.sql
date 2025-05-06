-- name: CreateTaskWithID :one
INSERT INTO tasks (id, data)
VALUES (?, ?)
RETURNING *;

-- name: CreateTask :one
INSERT INTO tasks (data)
VALUES (?)
RETURNING *;

-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY id ASC;
