-- name: CreateTaskWithID :one
INSERT INTO tasks (fsm_id, id, event)
VALUES (?, ?, ?)
RETURNING *;

-- name: CreateTask :one
INSERT INTO tasks (fsm_id, event)
VALUES (?, ?)
RETURNING *;

-- name: ListTasks :many
SELECT * FROM tasks
WHERE fsm_id = ?
ORDER BY id ASC;
